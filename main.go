package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	secretKey = []byte("yourKeyJWT")
)

func main() {
	r := gin.Default()

	r.POST("/login", loginHandler)

	auth := r.Group("/auth")
	auth.Use(authMiddleware)
	auth.GET("/protected", protectedHandler)

	r.Run(":8080")
}

func loginHandler(c *gin.Context) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "userExample"
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to create token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Missing authentication token"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication failed"})
		c.Abort()
		return
	}

	if token.Valid {
		c.Next()
	} else {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
	}
}

func protectedHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Protected route"})
}
