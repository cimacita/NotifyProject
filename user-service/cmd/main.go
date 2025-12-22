package main

import (
	"NotifyProject/user-service/config"
	"NotifyProject/user-service/internal/db"
	"NotifyProject/user-service/internal/user"
	"NotifyProject/user-service/pkg/auth"
	"NotifyProject/user-service/pkg/kafka"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers)
	defer kafkaProducer.Close()

	pool := db.Connect(cfg.Db.Dsn)
	defer pool.Close()

	router := http.NewServeMux()

	userRepository := user.NewRepository(pool)

	userService := user.NewService(userRepository, kafkaProducer)

	jwtManager := auth.NewJWTManager(cfg.Auth.Secret)

	user.NewHandler(router, userService, jwtManager)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
