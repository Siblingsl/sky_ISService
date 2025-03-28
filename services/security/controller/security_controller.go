package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/security/dto"
	"sky_ISService/services/security/service"
	"sky_ISService/utils"
)

type SecurityController struct {
	service *service.SecurityService
}

func NewSecurityController(securityService *service.SecurityService) *SecurityController {
	return &SecurityController{
		service: securityService,
	}
}

func (c *SecurityController) SecurityControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	securityGroup := r.Group("/security")

	// 登陆
	securityGroup.POST("/admins/login", func(ctx *gin.Context) {
		var req dto.SecurityAdminLoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误")
			return
		}
		token, err := c.service.AdminLogin(ctx, req)
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		// 返回成功响应
		utils.Success(ctx, token)
	})

	// 验证码
	securityGroup.GET("/admins/code", func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			utils.Error(ctx, http.StatusBadRequest, "邮箱不能为空")
		}
		err := c.service.SendEmailCode(ctx, email)
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(ctx, 1)
	})

	securityGroup.GET("/admins/test", func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			utils.Error(ctx, http.StatusBadRequest, "邮箱不能为空")
		}
		data, err := c.service.Testxxxx(ctx, email)
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.Success(ctx, data)
	})
}
