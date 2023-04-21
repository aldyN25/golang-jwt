package main

import (
	jwt "golang-jwt/jwt/v2"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Recovery())
	authMiddleware, err := jwt.GinJWt()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())

	}

	//TOKEN
	router.POST("/request-token", authMiddleware.LoginHandler)
	router.GET("/refresh-token", authMiddleware.RefreshHandler)

	//ROOT
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Selamat Data di DIGITAL AGRI API",
		})
	})

}
