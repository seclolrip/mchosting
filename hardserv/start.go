package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

// Args:
// Java Version, Min Mem, Max Mem, Log4j workaround, any command args needed for modded worlds including multiple jar files for forge or spigot
// Garbage collection commands, optimizations, etc

func Start(javaversion string, username string) {
	uuidv7, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	homeDir := os.Getenv("HOME")
	fpath := homeDir + "/mchosting/hardserv/pendingbuilds/" + uuidv7.String() + "-" + username
	keygen := exec.Command("ssh-keygen", "-t", "ed25519", "-f", fpath, "-N", "")

	genErr := keygen.Run()
	if genErr != nil {
		fmt.Println("Error running key-gen: ", genErr)
		return
	}

	pubKeyContents, err := os.ReadFile(fpath + ".pub")
	if err != nil {
		fmt.Println("Error reading public key file:", err)
		return
	}

	dockerBuildCmd := exec.Command("sudo", "docker", "build", "--build-arg", "JAVA_VERSION="+javaversion, "--build-arg", "USERUSERNAME="+username, "--build-arg", "SSHKEY="+string(pubKeyContents), "-t", "testing", "-f", "./Dockerfile", ".")
	dockerBuildCmd.Dir = homeDir + "/mchosting/hardserv"

	dockerBuildErr := dockerBuildCmd.Run()
	if dockerBuildErr != nil {
		fmt.Println("Error running docker build: ", dockerBuildErr.Error())
		return
	}

	dockerRunCmd := exec.Command("sudo", "docker", "run", "-d", "testing")
	dockerRunCmd.Dir = homeDir + "/mchosting/hardserv"

	dockerRunErr := dockerRunCmd.Run()
	if dockerRunErr != nil {
		fmt.Println("Error running docker run: ", dockerRunErr.Error())
		return
	}
}
