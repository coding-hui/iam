// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/pkg/code"
)

type Response struct {
	// Success request is successful
	Success bool `json:"success"`

	// Code defines the business error code.
	Code int `json:"code"`

	// Msg contains the detail of this message.
	// This message is suitable to be exposed to external
	Msg string `json:"msg"`

	// Data return data object
	Data interface{} `json:"data,omitempty"`

	// Total total of page
	Total int64 `json:"total,omitempty"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    code,
		Msg:     msg,
		Data:    data,
	})
}

func Ok(c *gin.Context) {
	Result(code.ErrSuccess, nil, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(code.ErrSuccess, nil, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(code.ErrSuccess, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(code.ErrSuccess, data, message, c)
}

func OkWithPage(result interface{}, total int64, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    code.ErrSuccess,
		Msg:     "success",
		Data:    result,
		Total:   total,
	})
}

func OkWithPageDetailed(result interface{}, total int64, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    code.ErrSuccess,
		Msg:     message,
		Data:    result,
		Total:   total,
	})
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
		klog.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), Response{
			Success:   false,
			Code:      coder.Code(),
			Msg:       coder.String(),
			Reference: coder.Reference(),
		})

		return
	}

	Fail(c)
}
