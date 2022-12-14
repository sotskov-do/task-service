package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server serverConfig
}

type serverConfig struct {
	Port               string `envconfig:"PORT" default:"3000"`
	GRPCAuth           string `envconfig:"GRPC_AUTH" default:"localhost:4000"`
	GRPCAnalytic       string `envconfig:"GRPC_ANALYTIC" default:"localhost:5000"`
	Profiling          bool   `envconfig:"PROFILING" default:"false"`
	PgUrl              string `envconfig:"PG_URL" default:"postgres://postgres:1111@localhost:5432/mtsteta"`
	JsonDbFile         string `envconfig:"JSON_DB_FILE" default:"db.jsonl"`
	KafkaUrl           string `envconfig:"KAFKA_URL" default:"kafka:9092"`
	KafkaAnalyticTopic string `envconfig:"KAFKA_ANALYTIC_TOPIC" default:"analytic"`
	EmailWorkers       int    `envconfig:"EMAIL_WORKERS" default:"5"`
	EmailRateLimit     int    `envconfig:"EMAIL_RATE_LIMIT" default:"3"`
}

func New() (*Config, error) {
	var c Config

	err := envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
