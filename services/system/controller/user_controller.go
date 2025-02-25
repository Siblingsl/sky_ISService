package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/system/service"
)

type UserController struct {
	userService *service.UserService
}

// NewUserController 构造函数，依赖注入 UserService
func NewUserController(userService *service.UserService) *UserController {
	controller := &UserController{userService: userService}
	return controller
}

func (c *UserController) UserControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	userGroup := r.Group("/system/user")

	// 管理员用户
	userGroup.GET("/admin", func(ctx *gin.Context) {
		message, err := c.userService.AddUser()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, gin.H{"message": message})
	})
}
