package controllers

import (
	"net/http"
	"todolist/config"
	"todolist/response"
	"todolist/services"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (con AuthController) Login(c *gin.Context) {
	var input struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if !utils.BindAndValidate(c, &input) {
		return // 綁定或驗證失敗，已經回傳錯誤了，直接結束
	}

	service := services.NewAuthService(c.Request.Context())
	data, err := service.Login(config.DB, input.Account, input.Password)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, data)
}

func (con AuthController) Register(c *gin.Context) {
	var input struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binging:"required"`
	}

	if !utils.BindAndValidate(c, &input) {
		return
	}

	service := services.NewAuthService(c.Request.Context())
	user, err := service.Register(config.DB, input.Account, input.Password)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, user)
}
