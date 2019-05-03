package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var servName = flag.String("name", "server", "The server name shown on /")
var servAddr = flag.String("addr", ":8080", "The server address")

func init() {
	flag.Parse()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"server": *servName,
		})
	})

	fmt.Println("Server running on:", *servAddr)
	fmt.Println(r.Run(*servAddr))
}
