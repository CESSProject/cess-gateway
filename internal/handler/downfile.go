package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DownfileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
