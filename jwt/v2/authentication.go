package jwt

import (
	"context"
	"log"
	"net/http"
	"strings"

	"golang-jwt/config"

	"github.com/gin-gonic/gin"
)

type TokenMapping struct {
	UserId    string `json:"userId"`
	LoginAt   string `json:"loginAt"`
	ExpiredAt string `json:"expiredAt"`
}

func TokenAuthenticationWithRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		redisDBToken := config.RedisDBToken()

		authToken := c.Request.Header.Get("Authorization")
		log.Println("Auth Token : ", authToken)

		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "403 Forbidden",
				"desc":    "Tidak Memliki Akses",
			})

			return
		}

		token := strings.TrimPrefix(authToken, "Bearer ")
		log.Println("Token : ", token)

		// Check Lifetime Token
		// open redis connection
		ctx := context.Background()
		redisClient := config.InitRedisConnection(redisDBToken)

		checkToken, err := redisClient.Get(ctx, token).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid Token || Token Is Expired",
				"desc":    "Expired",
			})

			return
		}

		log.Println("CHECK TOKEN : ", checkToken)

		if checkToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid Token || Token Is Expired",
				"desc":    "Expired",
			})

			return
		}

		c.Next()
	}
}
