package dto

// OperationLogListRequest 获取操作日志列表请求结构体
type OperationLogListRequest struct {
	Username  string `json:"username" form:"username"`
	Ip        string `json:"ip" form:"ip"`
	Path      string `json:"path" form:"path"`
	Status    int    `json:"status" form:"status"`
	PageNum   int    `json:"page" form:"page"`
	PageSize  int    `json:"per_page" form:"per_page"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

// DeleteOperationLogRequest 批量删除操作日志请求结构体
type DeleteOperationLogRequest struct {
	OperationLogIds []int64 `json:"operation_log_ids" form:"operation_log_ids"`
}
