package jwt

import (
	"golang-jwt/config"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type User struct {
	UserName string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// ================GENERATE TOKEN=============
func GinJWt() (*jwt.GinJWTMiddleware, error) {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: config.AppJWT,
		Key:   []byte(config.KeyJWT),
		// Timeout: time.Hour,
		Timeout: 168 * time.Hour,
		//Timeout: 1 * time.Minute,
		MaxRefresh: 168 * time.Hour,
		// MaxRefresh: time.Hour,
		//MaxRefresh:  1 * time.Minute,
		IdentityKey: config.IdentityJWT,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					config.IdentityJWT: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[config.IdentityJWT].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			userID := loginVals.Username
			password := loginVals.Password

			statusAuth, _ := config.GetAuth(userID, password, "FALSE")

			if statusAuth == 200 {
				return &User{
					UserName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			statusAuth, _ := config.GetAuth(data.(*User).UserName, "", "TRUE")

			if statusAuth == 200 {
				return true
			}

			log.Println("---USER NOT MATCH JWT---")

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Println("---UNAUTORIZED JWT---")

			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
		return nil, err
	}

	return authMiddleware, nil
}
