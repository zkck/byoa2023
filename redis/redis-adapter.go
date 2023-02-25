package redis

import (
	"context"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdbTexts = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDR"),
	Password: "", // no password set
	DB:       0,  // use default DB
})

var rdbTokes = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDR"),
	Password: "", // no password set
	DB:       1,  // use default DB
})

/* func main() {
	fmt.Println(Store("123", "test hallo"))
	fmt.Println(Search("123", "test1"))
	fmt.Println(GetText("123"))
	fmt.Println(Delete("123"))
} */

var ctx = context.Background()

func Store(uuid, value string) error {
	err := rdbTexts.Set(ctx, uuid, value, 0).Err()
	if err != nil {
		return err
	}

	splitValue := strings.Split(value, " ")
	for _, element := range splitValue {
		err = rdbTokes.Set(ctx, element+"##!##"+uuid, "true", 0).Err()
	}
	if err != nil {
		return err
	}
	return nil
}

// match / found / error
func Search(uuid, term string) (bool, bool, error) {
	_, err := rdbTexts.Get(ctx, uuid).Result()
	if err != nil {
		if err == redis.Nil {
			return false, false, nil
		}
		return false, false, err
	}

	_, err = rdbTokes.Get(ctx, term+"##!##"+uuid).Result()

	if err != nil {
		if err == redis.Nil {
			return false, true, nil
		}
		return false, true, err
	}
	return true, true, nil
}

func GetText(uuid string) (string, bool, error) {
	val, err := rdbTexts.Get(ctx, uuid).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func Delete(uuid string) (bool, error) {
	_, _, err := GetText(uuid)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	_, err = rdbTexts.Del(ctx, uuid).Result()
	/*
		splitValue := strings.Split(value, " ")
		for _, element := range splitValue {
			err = rdbTokes.Del(ctx, element+"##!##"+uuid).Err()
		}*/
	if err != nil {
		/*if err == redis.Nil {
			return false, nil
		}*/
		return false, err
	}
	return true, nil
}
