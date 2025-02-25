package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/system/service"
)

type RoleController struct {
	roleService *service.RoleService
}

func NewRoleController(roleService *service.RoleService) *RoleController {
	controller := &RoleController{roleService: roleService}
	return controller
}

func (c *RoleController) RoleControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	roleGroup := r.Group("/system/role")

	// 管理员角色
	roleGroup.GET("admin", func(ctx *gin.Context) {
		message, err := c.roleService.AddRole()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, gin.H{"message": message})
	})
}
