package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Aderon1333/TPES/internal/config"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

func NewMongoClient(ctx context.Context, sc config.MongoConfig, logger *logfacade.LogFacade) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(sc.Url).SetServerAPIOptions(serverAPI))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func DisconnectMongoClient(ctx context.Context, client *mongo.Client) error {
	if err := client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
