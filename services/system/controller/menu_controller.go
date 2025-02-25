package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/system/service"
)

type MenuController struct {
	menuService *service.MenuService
}

func NewMenuController(menuService *service.MenuService) *MenuController {
	controller := &MenuController{menuService: menuService}
	return controller
}

func (c *MenuController) MenuControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	menuGroup := r.Group("/system/menu")

	// 管理员菜单
	menuGroup.GET("/admin", func(ctx *gin.Context) {
		message, err := c.menuService.AddMenu()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, gin.H{"message": message})
	})
}
