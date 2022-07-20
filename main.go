package main

import (
	"MoeTracker/handler"
	"MoeTracker/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	redis.InitRedis()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/announce", handler.ReceiveAnnounce)

	r.Run()
}
