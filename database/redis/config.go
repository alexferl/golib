package redis

import (
	"github.com/spf13/pflag"
)

type Config struct {
	Addr       string
	DB         int
	ClientName string
	Username   string
	Password   string
	URL        string
}

var DefaultConfig = &Config{
	Addr:       "localhost:6379",
	DB:         0,
	ClientName: "",
	Username:   "",
	Password:   "",
}

const (
	RedisAddr       = "redis-addr"
	RedisDB         = "redis-db"
	RedisClientName = "redis-client-name"
	RedisUsername   = "redis-username"
	RedisPassword   = "redis-password"
	RedisURL        = "redis-url"
)

func (c *Config) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Addr, RedisAddr, c.Addr, "Redis addr")
	fs.IntVar(&c.DB, RedisDB, c.DB, "Redis database")
	fs.StringVar(&c.ClientName, RedisClientName, c.ClientName, "Redis client name")
	fs.StringVar(&c.Username, RedisUsername, c.Username, "Redis username")
	fs.StringVar(&c.Password, RedisPassword, c.Password, "Redis password")
	fs.StringVar(&c.URL, RedisURL, c.URL, "Redis connection string")
}
