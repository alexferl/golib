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
	TLS        *TLS
}

type TLS struct {
	CertFile string
	KeyFile  string
	CAFile   string
}

var DefaultConfig = &Config{
	Addr:       "localhost:6379",
	DB:         0,
	ClientName: "",
	Username:   "",
	Password:   "",
	TLS: &TLS{
		CertFile: "",
		KeyFile:  "",
		CAFile:   "",
	},
}

const (
	RedisAddr        = "redis-addr"
	RedisDB          = "redis-db"
	RedisClientName  = "redis-client-name"
	RedisUsername    = "redis-username"
	RedisPassword    = "redis-password"
	RedisURL         = "redis-url"
	RedisTLSCertFile = "redis-tls-cert-file"
	RedisTLSKeyFile  = "redis-tls-key-file"
	RedisTLSCAFile   = "redis-tls-ca-file"
)

func (c *Config) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Addr, RedisAddr, c.Addr, "Redis addr")
	fs.IntVar(&c.DB, RedisDB, c.DB, "Redis database")
	fs.StringVar(&c.ClientName, RedisClientName, c.ClientName, "Redis client name")
	fs.StringVar(&c.Username, RedisUsername, c.Username, "Redis username")
	fs.StringVar(&c.Password, RedisPassword, c.Password, "Redis password")
	fs.StringVar(&c.URL, RedisURL, c.URL, "Redis connection string")
	fs.StringVar(&c.TLS.CertFile, RedisTLSCertFile, c.TLS.CertFile, "Redis TLS certificate file")
	fs.StringVar(&c.TLS.KeyFile, RedisTLSKeyFile, c.TLS.KeyFile, "Redis TLS key file")
	fs.StringVar(&c.TLS.CAFile, RedisTLSCAFile, c.TLS.CAFile, "Redis TLS CA file")
}
