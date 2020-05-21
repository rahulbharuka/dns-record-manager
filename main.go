package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/dns-record-manager/external/route53"
	"github.com/rahulbharuka/dns-record-manager/logic"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// set release mode logging.
	gin.SetMode(gin.ReleaseMode)

	// create default Gin router
	router := gin.New()

	// init log middleware
	router.Use(gin.Logger())

	// init recovery middleware
	router.Use(gin.Recovery())

	// load HTML templates
	router.LoadHTMLGlob("./views/*.html")

	// initialize AWS Route53 service client.
	route53.Init()

	// get logic handler
	h := logic.GetHandler()

	// API handlers.
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/dns-records", h.ListDNSARecords)
	router.GET("/servers", h.ListServers)
	router.POST("/servers/:id/add", h.AddServer)
	router.POST("/servers/:id/remove", h.RemoveServer)

	// run app on the specified port
	router.Run(":" + port)
}
