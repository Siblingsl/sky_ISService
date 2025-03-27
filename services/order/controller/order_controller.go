package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/auth/service"
	"sky_ISService/utils"
)

// OrderController 处理身份验证相关的请求
type OrderController struct {
	service *service.AuthService
}

func NewOrderController(authService *service.AuthService) *OrderController {
	return &OrderController{
		service: authService,
	}
}

// OrderControllerRoutes 设置身份验证相关的路由
// @Summary 管理员相关操作
// @Description 提供管理员登录、验证码发送等功能
// @Tags Order
// @Accept json
// @Produce json
func (c *OrderController) OrderControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	orderGroup := r.Group("/order")

	orderGroup.GET("/getTest", func(ctx *gin.Context) {
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
