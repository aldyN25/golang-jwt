package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

const (
	AppJWT      = "digital-agri-api"
	IdentityJWT = "d1gitAL@4gR1"
	KeyJWT      = "ZGlnaXRhbC1hZ3JpOjIwMjEwNDIy"
)

func StringToInt(text string) (result int, status int) {
	result, _ = strconv.Atoi(text)

	return result, 200
}
func GetAuth(username string, password string, justUsername string) (status int, auth map[string]string) {
	userApi := os.Getenv("API_USER")
	secretApi := os.Getenv("API_SECRET")

	auth = map[string]string{}
	auth[userApi] = secretApi

	if auth[username] != "" { //check username
		if justUsername == "TRUE" {
			return 200, auth
		}

		if auth[username] == password { //check password
			return 200, auth
		} else {
			return 400, auth
		}

	} else {
		return 400, auth
	}
	return 200, auth
}

func RedisDBToken() (result string) {
	result = os.Getenv("REDIS_DB_TOKEN")
	return result
}

func RedisUrl() (name string, status int) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	stringUrl := []string{host, ":", port}
	result := strings.Join(stringUrl, "")
	return result, 200
}

func RedisPassword() (name string, status int) {
	result := os.Getenv("REDIS_PASSWORD")
	return result, 200
}

func RedisDB() (name int, status int) {
	dbRedis := os.Getenv("REDIS_DB")
	result, _ := StringToInt(dbRedis)
	return result, 200
}

func InitRedisConnection(redisDb string) (rdb *redis.Client) {
	redisUrl, _ := RedisUrl()
	redisPassword, _ := RedisPassword()
	redisDB, _ := RedisDB()
	if redisDb != "" {
		redisDB, _ = StringToInt(redisDb)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword, // no password set
		DB:       redisDB,       // use default DB
	})

	return rdb
}
