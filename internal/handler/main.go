package handler

import (
	"cess-gateway/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"}
	r.Use(cors.New(config))

	//
	r.PUT("/:filename", UpfileHandler)
	r.GET("/:filename", DownfileHandler)
	r.POST("/auth", GrantTokenHandler)
	r.GET("/files", FilelistHandler)
	r.DELETE("/:filename", DeletefileHandler)
	r.Run(":" + configs.Confile.ServicePort)
}
