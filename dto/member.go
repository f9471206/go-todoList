package dto

type MembeQuery struct {
	Page     int    `form:"page" example:"1" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" example:"10" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" example:"角色名稱"`
	Order    string `form:"order" example:"created_at desc"`
}

type MenberUpdate struct {
	RoleID int `json:"role_id" binding:"required"`
}
