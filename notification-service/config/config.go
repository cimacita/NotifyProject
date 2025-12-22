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
	Redis RedisConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

type KafkaConfig struct {
	Brokers string
	Topic   string
	GroupID string
}

type RedisConfig struct {
	Addr     string
	Password string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading notification-service/.env file")
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
			Topic:   os.Getenv("KAFKA_USER_EVENTS_TOPIC"),
			GroupID: os.Getenv("KAFKA_CONSUMER_GROUP_ID"),
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}
}
