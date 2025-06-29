package dto

type TodoTypeCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type TodoTypeQuery struct {
	Page     int    `form:"page" example:"1" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" example:"10" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" example:"任務類別"`
	Order    string `form:"order" example:"created_at desc"`
}

type TodoTypeUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
