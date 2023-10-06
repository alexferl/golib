package mongodb

import (
	"time"

	"github.com/spf13/pflag"
)

type Config struct {
	URI                      string
	AppName                  string
	Username                 string
	Password                 string
	ReplicaSet               string
	ServerSelectionTimeoutMs time.Duration
	ConnectTimeoutMs         time.Duration
	SocketTimeoutMs          time.Duration // query timeout
}

const (
	MongoDBURI                      = "mongodb-uri"
	MongoDBAppName                  = "mongodb-app-name"
	MongoDBUsername                 = "mongodb-username"
	MongoDBPassword                 = "mongodb-password"
	MongoDBReplicaSet               = "mongodb-replica-set"
	MongoDBServerSelectionTimeoutMs = "mongodb-server-selection-timeout-ms"
	MongoDBConnectTimeoutMs         = "mongodb-connect-timeout-ms"
	MongoDBSocketTimeoutMs          = "mongodb-socket-timeout-ms"
)

func (c *Config) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.URI, MongoDBURI, c.URI, "MongoDB URI")
	fs.StringVar(&c.AppName, MongoDBAppName, c.AppName, "MongoDB app name")
	fs.StringVar(&c.Username, MongoDBUsername, c.Username, "MongoDB username")
	fs.StringVar(&c.Password, MongoDBPassword, c.Password, "MongoDB password")
	fs.StringVar(&c.ReplicaSet, MongoDBReplicaSet, c.ReplicaSet, "MongoDB replica set")
	fs.DurationVar(&c.ServerSelectionTimeoutMs, MongoDBServerSelectionTimeoutMs,
		c.ServerSelectionTimeoutMs, "MongoDB server selection timeout ms")
	fs.DurationVar(&c.ConnectTimeoutMs, MongoDBConnectTimeoutMs, c.ConnectTimeoutMs,
		"MongoDB connect timeout ms")
	fs.DurationVar(&c.SocketTimeoutMs, MongoDBSocketTimeoutMs, c.SocketTimeoutMs,
		"MongoDB socket timeout ms")
}
