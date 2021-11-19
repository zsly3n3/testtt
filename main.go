package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	addr := fmt.Sprintf(`:%d`, 8180)
	err:=r.Run(addr)
	if err!=nil{
		log.Fatal(`run_err:`,err.Error())
	}
}