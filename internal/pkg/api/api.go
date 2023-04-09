package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	"github.com/wecoding/iam/internal/pkg/code"
)

type Response struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Msg contains the detail of this message.
	// This message is suitable to be exposed to external
	Msg string `json:"msg"`

	Data interface{} `json:"data"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// PageResponse page response
type PageResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Ok(c *gin.Context) {
	Result(code.ErrSuccess, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(code.ErrSuccess, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(code.ErrSuccess, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(code.ErrSuccess, data, message, c)
}

func OkWithPage(result interface{}, count int64, c *gin.Context) {
	Result(code.ErrSuccess, &PageResponse{
		List:  result,
		Total: count,
	}, "success", c)
}

func OkWithPageDetailed(result interface{}, count int64, message string, c *gin.Context) {
	Result(code.ErrSuccess, &PageResponse{
		List:  result,
		Total: count,
	}, message, c)
}

func Fail(c *gin.Context) {
	Result(code.ErrUnknown, map[string]interface{}{}, "failed", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(code.ErrUnknown, map[string]interface{}{}, message, c)
}

// FailWithErrCode write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func FailWithErrCode(err error, c *gin.Context) {
	if err != nil {
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), Response{
			Code:      coder.Code(),
			Msg:       coder.String(),
			Reference: coder.Reference(),
		})

		return
	}

	Fail(c)
}
