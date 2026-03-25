package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := rdb.Ping(ctx).Result()
	MustNot(err)
	fmt.Println(res)

	status, err := rdb.Set(ctx, "prem", "modha", 0).Result()
	MustNot(err)

	fmt.Println(status)

	getResult, err := rdb.Get(ctx, "aprem").Result()
	if err != nil {
		fmt.Print("couldnt find the key\n")
	} else {
		println(getResult)
	}
}

func MustNot(error error) {
	if error != nil {
		panic(error)
	}
}
