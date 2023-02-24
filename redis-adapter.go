package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func main() {
	fmt.Println(store("123", "test hallo"))
	fmt.Println(search("123", "test"))
	fmt.Println(getText("123"))
	fmt.Println(delete("123"))
}

var ctx = context.Background()

func store(uuid, value string) (bool, error) {
	err := rdb.Set(ctx, uuid, value, 0).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func search(uuid, term string) (bool, error) {
	val, err := rdb.Get(ctx, uuid).Result()

	if err != nil {
		return false, err
	}
	return strings.Contains(val, term), nil
}

func getText(uuid string) (string, error) {
	val, err := rdb.Get(ctx, uuid).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func delete(uuid string) (bool, error) {
	_, err := rdb.Del(ctx, uuid).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}