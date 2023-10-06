package mongodb

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := viper.GetString(MongoDBURI)
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	opts := options.Client()
	opts.ApplyURI(uri)

	appName := viper.GetString(MongoDBAppName)
	if appName != "" {
		opts.SetAppName(appName)
	}

	serverSelectionTimeout := viper.GetDuration(MongoDBServerSelectionTimeoutMs)
	if serverSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(serverSelectionTimeout)
	}

	connectTimeout := viper.GetDuration(MongoDBConnectTimeoutMs)
	if connectTimeout > 0 {
		opts.SetConnectTimeout(connectTimeout)
	}

	socketTimeout := viper.GetDuration(MongoDBSocketTimeoutMs)
	if socketTimeout > 0 {
		opts.SetSocketTimeout(socketTimeout)
	}

	username := viper.GetString(MongoDBUsername)
	password := viper.GetString(MongoDBPassword)
	if username != "" {
		opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	replSet := viper.GetString(MongoDBReplicaSet)
	if replSet != "" {
		opts.SetReplicaSet(replSet)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateIndexes(ctx context.Context, db *mongo.Database, indexes map[string][]mongo.IndexModel) error {
	for collection, models := range indexes {
		_, err := db.Collection(collection).Indexes().CreateMany(ctx, models)
		if err != nil {
			return err
		}
	}

	return nil
}
