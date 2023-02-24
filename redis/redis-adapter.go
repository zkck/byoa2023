package redis

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

/* func main() {
	fmt.Println(store("123", "test hallo"))
	fmt.Println(search("123", "test"))
	fmt.Println(getText("123"))
	fmt.Println(delete("123"))
} */

var ctx = context.Background()

func Store(uuid, value string) error {
	err := rdb.Set(ctx, uuid, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func Search(uuid, term string) (bool, error) {
	val, err := rdb.Get(ctx, uuid).Result()

	if err != nil {
		return false, err
	}
	return strings.Contains(val, term), nil
}

func GetText(uuid string) (string, error) {
	val, err := rdb.Get(ctx, uuid).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func Delete(uuid string) error {
	_, err := rdb.Del(ctx, uuid).Result()
	if err != nil {
		return err
	}
	return nil
}