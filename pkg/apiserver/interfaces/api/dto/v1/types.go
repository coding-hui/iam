package v1

// EmptyResponse empty response, it will used for delete api
type EmptyResponse struct{}

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page"`           // 页码
	PageSize int    `json:"page_size" form:"page_size"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`     //关键字
}
