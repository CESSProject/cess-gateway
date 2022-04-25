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
	config.AllowMethods = []string{"GET", "POST", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"}
	r.Use(cors.New(config))

	//
	r.PUT("/:filename", UpfileHandler)
	r.GET("/:filename", DownfileHandler)
	r.POST("/user/grant", GrantTokenHandler)
	r.GET("/file/list", FilelistHandler)
	r.GET("/user/state", UserStateHandler)
	r.POST("/file/delete", DeletefileHandler)
	r.GET("/space/price", QueryPriceHandler)
	r.Run(":" + configs.Confile.ServicePort)
}
