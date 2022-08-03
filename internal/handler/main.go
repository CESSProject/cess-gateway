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
	config.AddAllowHeaders("Authorization", "*")
	r.Use(cors.New(config))

	// handler
	r.PUT("/:filename", UpfileHandler)
	r.GET("/:fid", DownfileHandler)
	r.POST("/auth", GrantTokenHandler)
	r.GET("/files", FilelistHandler)
	r.GET("/state/:fid", FilestateHandler)
	r.DELETE("/:fid", DeletefileHandler)

	// run
	r.Run(":" + configs.C.ServicePort)
}
