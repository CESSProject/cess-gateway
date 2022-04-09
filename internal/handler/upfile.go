package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpfileHandler(c *gin.Context) {
	fmt.Println("upload")
	c.JSON(http.StatusOK, nil)
}
