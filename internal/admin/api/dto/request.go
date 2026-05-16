package dto

type ApiListRequest struct {
	Method    string `json:"method" form:"method"`
	Path      string `json:"path" form:"path"`
	Category  string `json:"category" form:"category"`
	Creator   string `json:"creator" form:"creator"`
	PageNum   int64  `json:"page" form:"page"`
	PageSize  int64  `json:"per_page" form:"per_page"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}
type CreateApiRequest struct {
	Method   string `json:"method" binding:"required,min=1,max=20"`
	Path     string `json:"path" binding:"required,min=1,max=100"`
	Category string `json:"category" binding:"required,min=1,max=50"`
	Desc     string `json:"desc" binding:"min=0,max=100"`
}
type UpdateApiRequest struct {
	Method   string `json:"method" binding:"min=1,max=20"`
	Path     string `json:"path" binding:"min=1,max=100"`
	Category string `json:"category" binding:"min=1,max=50"`
	Desc     string `json:"desc" binding:"min=0,max=100"`
}
type DeleteApiRequest struct {
	ApiIds []int64 `json:"api_ids"`
}
