package main

import (
	"NotifyProject/notification-service/config"
	"NotifyProject/notification-service/internal/db"
	"NotifyProject/notification-service/internal/events"
	"NotifyProject/notification-service/internal/notification"
	"NotifyProject/notification-service/internal/userShadow"
	"NotifyProject/notification-service/pkg/auth"
	"NotifyProject/notification-service/pkg/cache"
	k "NotifyProject/notification-service/pkg/kafka"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	pool := db.Connect(cfg.Db.Dsn)
	defer pool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	router := http.NewServeMux()

	notifRepository := notification.NewRepository(pool)
	userShadowRepository := userShadow.NewRepository(pool)

	c := cache.NewRedisCache[[]notification.Notification](redisClient)
	notifCache := notification.NewNotifCache(c)

	notifService := notification.NewService(notifRepository, notifCache, userShadowRepository)

	jwtManager := auth.NewJWTManager(cfg.Auth.Secret)

	notification.NewHandler(router, notifService, jwtManager)

	eventHandler := events.NewUserEventHandler(userShadowRepository, notifService)

	consumer := k.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topic, eventHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	consumer.Start()
	defer consumer.Close()

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
