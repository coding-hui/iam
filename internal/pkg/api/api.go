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

// Response Http API common response.
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

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// PageInfo Http API common page info.
type PageInfo struct {
	// List all records
	List interface{} `json:"list"`
	// Total all count
	Total int64 `json:"total"`
}

// Result build result info.
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    code,
		Msg:     msg,
		Data:    data,
	})
}

// PageResult build page result info.
func PageResult(code int, result interface{}, total int64, msg string, c *gin.Context) {
	Result(code, PageInfo{
		List:  result,
		Total: total,
	}, msg, c)
}

// Ok return success.
func Ok(c *gin.Context) {
	Result(code.ErrSuccess, nil, "success", c)
}

// OkWithMessage return success with message.
func OkWithMessage(message string, c *gin.Context) {
	Result(code.ErrSuccess, nil, message, c)
}

// OkWithData return success with data.
func OkWithData(data interface{}, c *gin.Context) {
	Result(code.ErrSuccess, data, "success", c)
}

// OkWithDetailed return success with data and message.
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(code.ErrSuccess, data, message, c)
}

// OkWithPage return success with page.
func OkWithPage(result interface{}, total int64, c *gin.Context) {
	PageResult(code.ErrSuccess, result, total, "success", c)
}

// OkWithPageDetailed return success with page.
func OkWithPageDetailed(result interface{}, total int64, message string, c *gin.Context) {
	PageResult(code.ErrSuccess, result, total, message, c)
}

// Fail return fail.
func Fail(c *gin.Context) {
	Result(code.ErrUnknown, map[string]interface{}{}, "failed", c)
}

// FailWithMessage return fail with message.
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
