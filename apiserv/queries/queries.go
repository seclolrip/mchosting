package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PendingServer struct {
	UserId                string    `bson:"userId"`
	UUID                  uuid.UUID `bson:"uuid"`
	ServerUUID            uuid.UUID `bson:"serverUUID"`
	ModsUploaded          bool      `bson:"modsUploaded"`
	ModpacksUploaded      bool      `bson:"modpacksUploaded"`
	ResourcepacksUploaded bool      `bson:"resourcepacksUploaded"`
}

type InternalServer struct {
	UUID   uuid.UUID `bson:"uuid"`
	PrivIP string    `bson:"privIP"`
}

func GetPendingServer(db *mongo.Database, userId string) (*PendingServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serverRaw := db.Collection("pendingServers").FindOne(ctx, bson.M{"userId": userId})
	if serverRaw.Err() != nil {
		return nil, serverRaw.Err()
	}

	var serverData PendingServer
	if err := serverRaw.Decode(&serverData); err != nil {
		return nil, serverRaw.Err()
	}

	return &serverData, nil
}

func GetInternalServer(db *mongo.Database, serverUUID uuid.UUID) (*InternalServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serverRaw := db.Collection("internalServers").FindOne(ctx, bson.M{"uuid": serverUUID})
	if serverRaw.Err() != nil {
		return nil, serverRaw.Err()
	}

	var serverData InternalServer
	if err := serverRaw.Decode(&serverData); err != nil {
		return nil, serverRaw.Err()
	}

	return &serverData, nil
}
