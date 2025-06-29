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

type MenberController struct{}

// @Summary 取得 Member 列表
// @Description 查詢 Member 清單，支援關鍵字與排序
// @Tags Member
// @Accept json
// @Produce json
// @Param page query int false "頁碼（預設 1）"
// @Param page_size query int false "每頁筆數（預設 10）"
// @Param keyword query string false "關鍵字搜尋"
// @Param order query string false "排序欄位與方式，如 created_at desc"
// @Security BearerAuth
// @Router /api/member [get]
func (con MenberController) Index(c *gin.Context) {
	var query dto.MembeQuery
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

	repo := repositories.NewAuthRepository()
	service := services.NewMemberService(c.Request.Context(), repo)

	result, err := service.Index(config.DB, query.Keyword, query.Page, query.PageSize, orders)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPagination(c, result)
}

// Show Member
// @Summary 取得單一 Member
// @Description 根據 ID 取得 Member 詳細資料
// @Tags Member
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {object} models.User "成功回傳 User"
// @Security BearerAuth
// @Router /api/member/{id} [get]
func (con MenberController) Show(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	repo := repositories.NewAuthRepository()
	service := services.NewMemberService(c.Request.Context(), repo)

	result, err := service.Show(config.DB, id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// Edit Member
// @Summary 修改 Member
// @Description 根據 ID 修改 Member 權限
// @Tags Member
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Param input body dto.MenberUpdate true "要更新的 Member 資料"
// @Success 200 {object} models.User "成功回傳更新後的 User"
// @Security BearerAuth
// @Router /api/member/{id} [put]
func (con MenberController) Edit(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input dto.MenberUpdate

	if !utils.BindAndValidate(c, &input) {
		return
	}

	repo := repositories.NewAuthRepository()
	service := services.NewMemberService(c.Request.Context(), repo)

	result, err := service.Edit(config.DB, id, input.RoleID)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)

}
