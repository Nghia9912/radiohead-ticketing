package main

import (
	"log"
	"net/http"
	"os"

	goredis "github.com/redis/go-redis/v9"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/NghiaHoang/radiohead-ticketing/internal/delivery/http/handler"
	"github.com/NghiaHoang/radiohead-ticketing/internal/repository/redis"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: redisAddr,
	})

	queueRepo := redis.NewQueueRedisRepo(redisClient)
	queueHandler := handler.NewQueueHandler(queueRepo)

	mux := http.NewServeMux()
	
	// Business API endpoints
	mux.HandleFunc("/api/v1/queue/enqueue", queueHandler.Enqueue)
	
	// Expose metrics endpoint for Prometheus to pull data
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
