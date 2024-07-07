package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Args:
// Java Version, Min Mem, Max Mem, Log4j workaround, any command args needed for modded worlds including multiple jar files for forge or spigot
// Garbage collection commands, optimizations, etc

type PendingServer struct {
	ObjectId primitive.ObjectID `bson:"_id"`
	JAVAV    string             `bson:"javav"`
	SSHKEY   string             `bson:"sshKey"`
}

const PENDINGBUILDS string = "pendingBuilds"

func GetDatabase() (*mongo.Database, error) {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connectString := os.Getenv("LIVECONNECT")
	if connectString == "" {
		panic("No Connect String Selected!")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := os.Getenv("DATABASE")
	if db == "" {
		panic("No DB Selected!")
	}

	return client.Database(db), nil
}

func CheckPendingBuilds(db *mongo.Database) {
	pendingBuildsTable := db.Collection(PENDINGBUILDS)

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer timeoutCancel()

	var findDocument PendingServer
	findDBRes := pendingBuildsTable.FindOneAndUpdate(timeoutCtx, bson.M{"status": "new"}, bson.M{"$set": bson.M{"status": "pending"}})
	if findDBRes.Err() != nil && !errors.Is(findDBRes.Err(), mongo.ErrNoDocuments) {
		if errors.Is(findDBRes.Err(), mongo.ErrClientDisconnected) || mongo.IsNetworkError(findDBRes.Err()) {
			panic("Find failed with fatal error " + findDBRes.Err().Error())
		}
		return
	}

	decodeErr := findDBRes.Decode(&findDocument)
	if decodeErr != nil {
		panic("Failed to decode MongoDB Pending Server " + decodeErr.Error())
	}

	pubKey, privKey, buildAndRunErr := BuildAndRun(db, findDocument.JAVAV)
	if buildAndRunErr != nil {
		timeoutCtx, updateCancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
		defer updateCancel()

		update, updateErr := pendingBuildsTable.UpdateOne(timeoutCtx, bson.M{"_id": findDocument.ObjectId}, bson.M{"$set": bson.M{"status": "errored", "error": buildAndRunErr.Error()}})
		if updateErr != nil {
			panic("Failed to update MongoDB Failed Pending Build " + updateErr.Error())
		} else if update.MatchedCount == 0 {
			panic("CONFLICT: Failed Docker Build but didn't update MongoDB")
		}
	} else {
		timeoutCtx, updateCancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
		defer updateCancel()

		update, updateErr := pendingBuildsTable.UpdateOne(timeoutCtx, bson.M{"_id": findDocument.ObjectId}, bson.M{"$set": bson.M{"status": "completed", "pubKey": pubKey, "privKey": privKey}})
		if updateErr != nil {
			panic("Failed to update MongoDB Created Pending Build " + updateErr.Error())
		} else if update.MatchedCount == 0 {
			panic("CONFLICT: Created Docker Build but didn't update MongoDB")
		}
	}
}

// Returns: Pub Key, Priv Key, Error
func BuildAndRun(db *mongo.Database, javaversion string) (string, string, error) {
	uuidv7, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	homeDir := os.Getenv("HOME")
	fpath := homeDir + "/builds/" + uuidv7.String()

	mkdir := exec.Command("mkdir", fpath)
	mkdirErr := mkdir.Run()
	if mkdirErr != nil {
		return "", "", fmt.Errorf("error running mkdir: %v", mkdirErr.Error())
	}

	keygen := exec.Command("ssh-keygen", "-t", "ed25519", "-f", fpath+"/sshkey", "-N", "")
	defer clearTemp(uuidv7.String())

	genErr := keygen.Run()
	if genErr != nil {
		return "", "", fmt.Errorf("error running key-gen: %v", genErr.Error())
	}

	pubKeyContents, err := os.ReadFile(fpath + "/sshkey.pub")
	if err != nil {
		return "", "", fmt.Errorf("error reading public key file: %v", err.Error())
	}
	privKeyContents, privErr := os.ReadFile(fpath + "/sshkey")
	if privErr != nil {
		return "", "", fmt.Errorf("error reading private key file: %v", privErr.Error())
	}

	dockerBuildCmd := exec.Command("sudo", "docker", "build", "--build-arg", "JAVA_VERSION="+javaversion, "--build-arg", "MCFILES="+fpath+"/mcfiles", "--build-arg", "SSHKEY="+string(pubKeyContents), "-t", uuidv7.String(), "-f", homeDir+"/mchosting/hardserv/docker/Dockerfile", ".")

	dockerBuildErr := dockerBuildCmd.Run()
	if dockerBuildErr != nil {
		return "", "", fmt.Errorf("error running docker build: %v", dockerBuildErr.Error())
	}

	dockerRunCmd := exec.Command("sudo", "docker", "run", "-d", uuidv7.String())

	dockerRunErr := dockerRunCmd.Run()
	if dockerRunErr != nil {
		return "", "", fmt.Errorf("error running docker run: %v", dockerRunErr.Error())
	}

	return strings.Replace(string(pubKeyContents), "\n", "", -1), strings.Replace(string(privKeyContents), "\n", "", -1), nil
}

func clearTemp(uuidv7 string) {
	homeDir := os.Getenv("HOME")
	fpath := homeDir + "/builds/" + uuidv7
	os.Remove(fpath)
}

func main() {
	var osSigChan chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGTERM)

	mongodb, dbErr := GetDatabase()
	if dbErr != nil {
		panic("Failed to connect to MongoDB " + dbErr.Error())
	}

	cronSch, cronSchErr := gocron.NewScheduler(gocron.WithLimitConcurrentJobs(5, gocron.LimitModeReschedule))
	if cronSchErr != nil {
		panic("Failed to create Cron Scheduler " + cronSchErr.Error())
	}

	_, addJobErr := cronSch.NewJob(
		gocron.CronJob("*/5 * * * * *", true),
		gocron.NewTask(CheckPendingBuilds, mongodb),
	)
	if addJobErr != nil {
		panic(addJobErr.Error())
	}

	cronSch.Start()

	<-osSigChan
	err := cronSch.Shutdown()
	if err != nil {
		panic("Failed to shutdown cron " + err.Error())
	}
}
