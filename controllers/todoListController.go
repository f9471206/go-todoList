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

type TodoListController struct{}

// Create TodoList
// @Summary 新增 TodoList
// @Description 建議一個新的 todoList 項目
// @Tags TodoList
// @Accept json
// @Produce json
// @Param input body dto.TodoListCreateRequest true "建立 TodoList 所需資料"
// @Success 200 {object} models.TodoList "建立成功回傳的 TodoList 資料"
// @Security BearerAuth
// @Router /api/todo/list [post]
func (ctl *TodoListController) Create(c *gin.Context) {
	var input dto.TodoListCreateRequest

	if !utils.BindAndValidate(c, &input) {
		return // 綁定或驗證失敗，已經回傳錯誤了，直接結束
	}

	repo := repositories.NewTodoListRepository()
	service := services.NewTodoListService(c.Request.Context(), repo)
	result, err := service.Create(config.DB, input.Name, input.TypeID) // 這裡多傳入 type_id
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// @Summary 取得 TodoList 列表
// @Description 查詢 TodoList 清單，支援關鍵字與排序
// @Tags TodoList
// @Accept json
// @Produce json
// @Param page query int false "頁碼（預設 1）"
// @Param page_size query int false "每頁筆數（預設 10）"
// @Param keyword query string false "關鍵字搜尋"
// @Param order query string false "排序欄位與方式，如 created_at desc"
// @Security BearerAuth
// @Router /api/todo/list [get]
func (ctl *TodoListController) Index(c *gin.Context) {
	// 取得查詢參數
	var query dto.TodoListQuery

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

	repo := repositories.NewTodoListRepository()
	service := services.NewTodoListService(c.Request.Context(), repo)

	result, err := service.Index(config.DB, query.Keyword, query.Page, query.PageSize, orders)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPagination(c, result)
}

// Show TodoList
// @Summary 取得單一 TodoList
// @Description 根據 ID 取得 TodoList 詳細資料
// @Tags TodoList
// @Accept json
// @Produce json
// @Param id path int true "TodoList ID"
// @Success 200 {object} models.TodoList "成功回傳 TodoList"
// @Security BearerAuth
// @Router /api/todo/list/{id} [get]
func (ctl *TodoListController) Show(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	repo := repositories.NewTodoListRepository()
	service := services.NewTodoListService(c.Request.Context(), repo)
	result, err := service.Show(config.DB, id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// Edit TodoList
// @Summary 修改 TodoList
// @Description 根據 ID 修改 TodoList 名稱
// @Tags TodoList
// @Accept json
// @Produce json
// @Param id path int true "TodoList ID"
// @Param input body dto.TodeListUpdateRequest true "要更新的 TodoList 資料"
// @Success 200 {object} models.TodoList "成功回傳更新後的 TodoList"
// @Security BearerAuth
// @Router /api/todo/list/{id} [put]
func (ctl *TodoListController) Edit(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input dto.TodeListUpdateRequest
	if !utils.BindAndValidate(c, &input) {
		return // 綁定或驗證失敗，已經回傳錯誤了，直接結束
	}

	repo := repositories.NewTodoListRepository()
	service := services.NewTodoListService(c.Request.Context(), repo)
	result, err := service.Edit(config.DB, id, input.Name, input.TypeID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// Delete TodoList
// @Summary 刪除 TodoList
// @Description 根據 ID 刪除指定的 TodoList
// @Tags TodoList
// @Accept json
// @Produce json
// @Param id path int true "TodoList ID"
// @Success 200 {object} models.TodoList "成功回傳被刪除的 TodoList"
// @Security BearerAuth
// @Router /api/todo/list/{id} [delete]
func (ctl *TodoListController) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	repo := repositories.NewTodoListRepository()
	service := services.NewTodoListService(c.Request.Context(), repo)
	result, err := service.Delete(config.DB, id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)

}
