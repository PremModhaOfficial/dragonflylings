package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	clent := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		MaxRetries:  1,
		PoolTimeout: 10 * time.Second,
	})
	clent.Ping(context.Background())

	a, err := clent.SetArgs(context.Background(), "prem", "orem", redis.SetArgs{
		Mode: "NX",
	}).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}
