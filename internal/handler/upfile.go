package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpfileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
