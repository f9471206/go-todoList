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

type TodoTypeController struct{}

// Create TodoType
// @Summary 新增 TodoType
// @Description 建立一個新的 TodoType 項目
// @Tags TodoTypes
// @Accept json
// @Produce json
// @Param input body dto.TodoTypeCreateRequest true "建立 TodoType 所需資料"
// @Success 200 {object} models.TodoTypes "建立成功回傳的 TodoType 資料"
// @Security BearerAuth
// @Router /api/todo/type [post]
func (ctl *TodoTypeController) Create(c *gin.Context) {
	var input dto.TodoTypeCreateRequest

	if !utils.BindAndValidate(c, &input) {
		return // 綁定或驗證失敗，已經回傳錯誤了，直接結束
	}

	// 建立 repo 實例（泛型已指定型別）
	repo := repositories.NewTodoTypeRepository()
	service := services.NewTodoTypeService(c.Request.Context(), repo)
	result, err := service.Create(config.DB, input.Name)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// @Summary 取得 TodoType 列表
// @Description 查詢 TodoType 清單，支援關鍵字與排序
// @Tags TodoTypes
// @Accept json
// @Produce json
// @Param page query int false "頁碼（預設 1）"
// @Param page_size query int false "每頁筆數（預設 10）"
// @Param keyword query string false "關鍵字搜尋"
// @Param order query string false "排序欄位與方式，如 created_at desc"
// @Security BearerAuth
// @Router /api/todo/type [get]
func (ctl *TodoTypeController) Index(c *gin.Context) {
	var query dto.TodoTypeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "無效的查詢參數")
		return
	}

	// 補預設值
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	orders := utils.ParseOrders(query.Order, utils.AllowedOrders, "created_at desc")

	repo := repositories.NewTodoTypeRepository()
	service := services.NewTodoTypeService(c.Request.Context(), repo)
	result, err := service.Index(config.DB, query.Keyword, query.Page, query.PageSize, orders)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPagination(c, result)
}

// Show TodoType
// @Summary 取得單一 TodoType
// @Description 根據 ID 取得 TodoType 詳細資料
// @Tags TodoTypes
// @Accept json
// @Produce json
// @Param id path int true "TodoType ID"
// @Success 200 {object} models.TodoTypes "成功回傳 TodoType"
// @Security BearerAuth
// @Router /api/todo/type/{id} [get]
func (ctl *TodoTypeController) Show(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	repo := repositories.NewTodoTypeRepository()
	service := services.NewTodoTypeService(c.Request.Context(), repo)
	result, err := service.Show(config.DB, id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// Edit TodoType
// @Summary 修改 TodoType
// @Description 根據 ID 修改 TodoType 名稱
// @Tags TodoTypes
// @Accept json
// @Produce json
// @Param id path int true "TodoType ID"
// @Param input body dto.TodoTypeUpdateRequest true "要更新的 TodoType 資料"
// @Success 200 {object} models.TodoTypes "成功回傳更新後的 TodoType"
// @Security BearerAuth
// @Router /api/todo/type/{id} [put]
func (ctl *TodoTypeController) Edit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input dto.TodoTypeUpdateRequest
	if !utils.BindAndValidate(c, &input) {
		return
	}

	repo := repositories.NewTodoTypeRepository()
	service := services.NewTodoTypeService(c.Request.Context(), repo)
	result, err := service.Edit(config.DB, id, input.Name)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// Delete TodoType
// @Summary 刪除 TodoType
// @Description 根據 ID 刪除指定的 TodoType
// @Tags TodoTypes
// @Accept json
// @Produce json
// @Param id path int true "TodoType ID"
// @Success 200 {object} models.TodoTypes "成功回傳被刪除的 TodoType"
// @Security BearerAuth
// @Router /api/todo/type/{id} [delete]
func (ctl *TodoTypeController) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的 ID")
		return
	}

	repo := repositories.NewTodoTypeRepository()
	service := services.NewTodoTypeService(c.Request.Context(), repo)
	result, err := service.Delete(config.DB, id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
