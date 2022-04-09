package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DownfileHandler(c *gin.Context) {
	fmt.Println("download")
	c.JSON(http.StatusOK, nil)
}
