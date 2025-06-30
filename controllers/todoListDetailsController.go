package controllers

import (
	"net/http"
	"strconv"
	"todolist/config"
	"todolist/dto"
	"todolist/repositories"
	"todolist/response"
	"todolist/services"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

type TodoListDetailsController struct{}

// Create TodoListDetails
// @summary 新增 todoListDetails
// @Description 建立一個新的 todoListDetails 項目
// @Tags TodoListDetails
// @Accept json
// @Produce json
// @Parm input body dto.TodoListDetailsCreateRequest true "建立 todoListDetails 所需資料"
// @Success 200 {object} models.TodoListDetails "建立成功回傳的 TodoListDetails 資料"
// @Security BearerAuth
// @Router /api/todo/details [post]
func (ctl *TodoListDetailsController) Create(c *gin.Context) {
	var input dto.TodoListDetailsCreateRequest

	if !utils.BindAndValidate(c, &input) {
		return // 綁定或驗證失敗，已經回傳錯誤了，直接結束
	}

	repo := repositories.NewTodoListDetailsRepository()
	service := services.NewTodoListDetailsService(c.Request.Context(), repo)
	result, err := service.Create(config.DB, input.TodoListID, input.Name, input.Detail, input.IDs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// Edit TodoListDetails
// @Summary 修改 TodoListDetails
// @Description 根據 ID 修改 TodoListDetails 名稱
// @Tags TodoListDetails
// @Accept json
// @Produce json
// @Param id path int true "TodoListDetails ID"
// @Param input body dto.TodoTypeUpdateRequest true "要更新的 TodoListDetails 資料"
// @Success 200 {object} models.TodoListDetails "成功回傳更新後的 TodoType"
// @Security BearerAuth
// @Router /api/todo/details/{id} [put]
func (ctl *TodoListDetailsController) Edit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input dto.TodoListDetailsUpdateRequest
	if !utils.BindAndValidate(c, &input) {
		return // 綁定或驗證失敗，已經回傳錯誤了，直接結束
	}

	repo := repositories.NewTodoListDetailsRepository()
	service := services.NewTodoListDetailsService(c.Request.Context(), repo)
	result, err := service.Edit(config.DB, id, input.Name, input.Detail)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

func (ctl *TodoListDetailsController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	repo := repositories.NewTodoListDetailsRepository()
	service := services.NewTodoListDetailsService(c.Request.Context(), repo)

	result, err := service.Delete(config.DB, id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
