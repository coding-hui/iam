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
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

const (
	ERROR   = "999"
	SUCCESS = "000"
)

func Result(code string, data interface{}, msg string, c *gin.Context) {
	// 开始时间
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

func OkWithPage(result interface{}, count int64, pageIndex, pageSize int, c *gin.Context) {
	Result(SUCCESS, &PageResponse{
		List:     result,
		Total:    count,
		PageSize: pageSize,
		Page:     pageIndex,
	}, "success", c)
}

func OkWithPageDetailed(result interface{}, count int64, pageIndex, pageSize int, message string, c *gin.Context) {
	Result(SUCCESS, &PageResponse{
		List:     result,
		Total:    count,
		PageSize: pageSize,
		Page:     pageIndex,
	}, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "failed", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
