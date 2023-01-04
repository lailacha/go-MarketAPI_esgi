package main

import (
	"log"
	//jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gin-gonic/gin"
)


func logWithToken(next gin.HandlerFunc) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		log.Println(token)
		next(c)
	}
	
}