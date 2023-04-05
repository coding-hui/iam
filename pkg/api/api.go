package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// PageResponse page response
type PageResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

const (
	ERROR   = "999"
	SUCCESS = "000"
)

func Result(code string, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func OkWithPage(result interface{}, count int64, c *gin.Context) {
	Result(SUCCESS, &PageResponse{
		List:  result,
		Total: count,
	}, "success", c)
}

func OkWithPageDetailed(result interface{}, count int64, message string, c *gin.Context) {
	Result(SUCCESS, &PageResponse{
		List:  result,
		Total: count,
	}, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "failed", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}
