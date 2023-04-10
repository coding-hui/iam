package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	defaultPageSize = "10"
	pageParam       = "offset"
	pageSizeParam   = "limit"
)

// ExtractPagingParams extract `page` and `pageSize` params from request
func ExtractPagingParams(c *gin.Context, minPageSize, maxPageSize int64) (int64, int64, error) {
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
	return page, pageSize, nil
}