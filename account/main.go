package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting server...")

	router := gin.Default()

	router.GET("/api/account", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H {
			"hello": "world",
		})
	})

	srv := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	srv.ListenAndServe()
}
