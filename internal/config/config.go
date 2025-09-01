package config

import (
	"github.com/joeshaw/envdecode"
)

var config *Config

type (
	Config struct {
		Server
		Database
		NATS
		Logger
		JWT
		Binance
	}
	// Server configurations
	Server struct {
		Port string `env:"SERVER_PORT,default=8080"`
		Env  string `env:"SERVER_ENV,default=local"`
	}
	// Database configurations
	Database struct {
		DSN string `env:"DATABASE_DSN,default=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	}
	NATS struct {
		NatsUsername       string `env:"NATS_USERNAME,default=admin"`
		NatsPassword       string `env:"NATS_PASSWORD,default=admin"`
		NatsClusterAddress string `env:"NATS_CLUSTER_ADDRESS,default=localhost"`
		NatsClusterPort    string `env:"NATS_CLUSTER_PORT,default=4222"`
		NatsClientID       string `env:"NATS_CLIENT_ID,default=nutriflex"`
		NatsSubscriberName string `env:"NATS_SUBSCRIBER_NAME,default=nutriflex-subscriber"`
		NatsStreamName     string `env:"NATS_STREAM_NAME,default=nutriflex-stream"`
		NatsDLQSubject     string `env:"NATS_DLQ_SUBJECT,default=nutriflex.dlq"`
		NatsFakeServerPort string `env:"NATS_FAKE_SERVER_PORT,default=1234"`
		// Domain events
		NatsDomainEventsPattern           string `env:"NATS_DOMAIN_EVENTS_PATTERN,default=endurance.events.*"`
		NatsMarketDataDomainEventsSubject string `env:"NATS_MARKET_DATA_DOMAIN_EVENTS_SUBJECT,default=endurance.events.market_data"`
	}
	Logger struct {
		Level          int64 `env:"LOG_LEVEL,default=4"`
		LogTestQueries bool  `env:"LOG_TEST_QUERIES"`
	}
	JWT struct {
		SigningKey string `env:"JWT_SIGNING_KEY,default=secret"`
	}
	Binance struct {
		BaseURL   string `env:"BINANCE_BASE_URL,default=https://api.binance.com"`
		APIKey    string `env:"BINANCE_API_KEY"`
		APISecret string `env:"BINANCE_API_SECRET"`
	}
)

func initCfg() {
	if config != nil {
		return
	}
	config = &Config{}
	if err := envdecode.Decode(config); err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	initCfg()
	return config
}
