package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db    DbConfig
	Auth  AuthConfig
	Kafka KafkaConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

type KafkaConfig struct {
	Brokers string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading user-service/.env file")
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
		Kafka: KafkaConfig{
			Brokers: os.Getenv("KAFKA_BROKERS"),
		},
	}
}
