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
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"}
	r.Use(cors.New(config))

	//
	r.POST("/file/upload", UpfileHandler)
	r.GET("/file/download", DownfileHandler)
	r.POST("/user/randoms", GenerateRandomkeyHandler)
	r.POST("/user/grant", GrantTokenHandler)
	//r.POST("/user/regrant", RegrantTokenHandler)
	r.GET("/file/list", FilelistHandler)
	r.GET("/user/state", UserStateHandler)
	r.POST("/file/delete", DeletefileHandler)
	r.GET("/space/price", QueryPriceHandler)
	r.Run(":" + configs.Confile.ServicePort)
}
