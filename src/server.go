package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/nats-io/nats.go"
)

var PARALLEL_BATCH, _ = strconv.Atoi(os.Getenv("PARALLEL_BATCH"))

func main() {
	server := getInitialServer()
	
	server.Use(serverEssentials())

	server.Use(getCustomRecoveryMiddleware())

	server.GET("/", func(c *gin.Context) {
		natsConn, err := nats.Connect(os.Getenv("NATS_SERVERS"))
		if err != nil {
			log.Fatalf("Could not connect to Nats Servers: %v", err, os.Getenv("NATS_SERVERS"))
		}
		defer natsConn.Close()
		uid, _ := uuid.NewV4()
		subject := uid.String()

		var wg sync.WaitGroup
		for i:=1; i<=PARALLEL_BATCH; i++ {
			wg.Add(1)
			go NatsLoadTest(natsConn, fmt.Sprintf("%s-%d", subject, i), &wg)
		}
		wg.Wait()

		c.JSON(200, gin.H{
			"Message": "OK !",
		})
	})

	server.Run(":9999")
}
