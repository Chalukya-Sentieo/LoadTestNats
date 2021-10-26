package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := getInitialServer()

	server.Use(getCustomRecoveryMiddleware())
	server.Use(getCORSMiddleware())

	server.GET("/", func(c *gin.Context) {
		testGetParam := c.Request.FormValue("testGetParam")
		if testGetParam == "BadRequest" {
			err := BadRequestException{}
			err.captureException("Testing a Bad Request and sentry integration")
		} else if testGetParam == "Unauthorized" {
			err := UnauthorizedRequestException{}
			err.captureException("Testing a Unauthorized request and sentry integration")
		} else if testGetParam == "Panic" {
			log.Panic("Throw exception using panic")
		} else if testGetParam == "NullPointer" {
			log.Println("Cause panic by accessing a null pointer")
			c.Request.GetBody()
		}
		c.JSON(http.StatusOK, gin.H{
			"Message": "OK !",
		})
	})

	server.Run(":80")
}
