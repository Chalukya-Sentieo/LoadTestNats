package main

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
)

func getInitialServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(gin.Logger())
	return server
}


/**
RequestID Middleware
*/

func serverEssentials() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uid, _ := uuid.NewV4()
        c.Header("X-Request-ID", uid.String())

        c.Next()
		log.Printf("Request Completed In : %v", time.Since(start).Truncate(time.Microsecond))
	}
}
