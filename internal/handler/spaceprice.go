package handler

import (
	"cess-httpservice/internal/chain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func QueryPriceHandler(c *gin.Context) {
	var resp = RespMsg{
		Code: http.StatusInternalServerError,
		Msg:  "",
	}
	price, err := queryPrice()
	if err != nil {
		resp.Msg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Code = http.StatusOK
	resp.Msg = "success"
	fmt.Println(price)
	if price == float64(0) {
		resp.Data = "space is sold out"
	} else {
		resp.Data = price
	}
	c.JSON(http.StatusOK, resp)
	return
}

/*
FindPrice means to get real-time price of storage space
*/
func queryPrice() (float64, error) {

	soldspace, err := chain.QuerySoldSpace()
	if err != nil {
		return 0, errors.Wrap(err, "QuerySoldSpace")
	}

	totalspace, err := chain.QueryTotalSpace()
	if err != nil {
		return 0, errors.Wrap(err, "QueryTotalSpace")
	}

	if soldspace == totalspace || totalspace < soldspace {
		return 0, nil
	}

	result := (1024 / float64((totalspace - soldspace))) * 1000

	return result, nil
}
