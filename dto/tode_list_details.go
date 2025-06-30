package dto

type TodoListDetailsCreateRequest struct {
	TodoListID int    `json:"to_do_list_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Detail     string `json:"detail" binding:"required"`
	IDs        []int  `json:"user_ids" binding:"required"`
}

type TodoListDetailsUpdateRequest struct {
	Name   string `json:"name" binding:"required"`
	Detail string `json:"detail" binding:"required"`
}
