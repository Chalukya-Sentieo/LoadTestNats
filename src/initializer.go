package main

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getInitialServer() *gin.Engine {
	server := gin.New()
	server.Use(gin.Logger())
	return server
}

func getCORSMiddleware() gin.HandlerFunc {
	// Get comma seperated list of origins from env and convert it array of ALLOWED ORIGINS
	originsAllowed := strings.Split(os.Getenv("ALLOWED_CORS_ORIGINS"), ",")

	// NOTE: In case a more complex CORS logic is required to authenticate origin,
	// there is a seperate field called AllowedOriginsFunc which will be a function that you write
	// and which can evaluate the ORIGIN based on a different logic and return a boolean true if it needs to be allowed
	return cors.New(cors.Config{
		AllowOrigins:     originsAllowed,
		AllowMethods:     []string{"PUT", "GET", "POST", "UPDATE", "OPTIONS", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "User-Agent", "Accept-Encoding", "Access-Control-Request-Headers", "Access-Control-Request-Method"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
