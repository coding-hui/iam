// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	defaultPageSize = "10"
	pageParam       = "current"
	pageSizeParam   = "pageSize"
	offsetParam     = "offset"
	limitParam      = "limit"
)

// ExtractPagingParams extract `page` and `pageSize` params from request
func ExtractPagingParams(c *gin.Context, minPageSize, maxPageSize int64) (int64, int64, error) {
	offsetStr := c.Query(offsetParam)
	limitStr := c.Query(limitParam)
	if offsetStr != "" && limitStr != "" {
		offset, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			return 0, 0, errors.Errorf("invalid offset %s: %v", offsetStr, err)
		}
		limit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			return 0, 0, errors.Errorf("invalid limit %s: %v", limitStr, err)
		}
		return offset, limit, nil
	}
	pageStr := c.Query(pageParam)
	pageSizeStr := c.Query(pageSizeParam)
	if pageStr == "" {
		pageStr = "0"
	}
	if pageSizeStr == "" {
		pageSizeStr = defaultPageSize
	}
	page64, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		return 0, 0, errors.Errorf("invalid page %s: %v", pageStr, err)
	}
	pageSize64, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		return 0, 0, errors.Errorf("invalid pageSize %s: %v", pageSizeStr, err)
	}
	page := page64
	pageSize := pageSize64
	if page < 0 {
		page = 0
	}
	if pageSize < minPageSize {
		pageSize = minPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return (page - 1) * pageSize64, pageSize, nil
}
