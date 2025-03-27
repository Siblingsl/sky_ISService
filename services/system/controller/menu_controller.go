package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/pkg/middleware"
	"sky_ISService/services/system/dto"
	"sky_ISService/services/system/service"
	"sky_ISService/utils"
	"strconv"
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
	menuGroup := r.Group("/system")

	// 添加菜单
	menuGroup.POST("/menu", func(ctx *gin.Context) {
		var req dto.CreateSkySystemMenuRequest
		// 1. 解析请求 JSON
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		// 2. 调用服务层创建菜单
		menu, err := c.menuService.CreateMenu(req)
		if err != nil {
			// 处理 "当前用户已存在" 错误
			if err.Error() == "当前菜单已存在" || err.Error() == "没有创建权限" {
				utils.Error(ctx, http.StatusBadRequest, err.Error()) // 400 用户请求错误
			} else {
				utils.Error(ctx, http.StatusInternalServerError, "创建用户失败: "+err.Error()) // 500 服务器错误
			}
			return
		}

		// 3. 返回成功响应
		utils.Success(ctx, menu)
	})

	// 获取菜单列表
	menuGroup.GET("/menu/list", func(ctx *gin.Context) {
		// 调用服务层获取所有菜单列表
		menus, err := c.menuService.GetMenuList()
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "获取菜单列表失败: "+err.Error()) // 500 服务器错误
			return
		}
		// 返回成功响应
		utils.Success(ctx, menus)
	})

	// 更新菜单
	menuGroup.PUT("/menu", func(ctx *gin.Context) {
		var req dto.UpdateSkySystemMenuRequest
		// 1. 解析请求 JSON
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error()) // 400 请求错误
			return
		}

		// 2. 调用服务层更新菜单
		menu, err := c.menuService.UpdateMenu(req)
		if err != nil {
			// 处理 "菜单不存在" 错误
			if err.Error() == "菜单不存在" {
				utils.Error(ctx, http.StatusNotFound, err.Error()) // 404 菜单不存在
			} else {
				utils.Error(ctx, http.StatusInternalServerError, "更新菜单失败: "+err.Error()) // 500 服务器错误
			}
			return
		}

		// 3. 返回成功响应
		utils.Success(ctx, menu)
	})

	// 删除菜单
	menuGroup.DELETE("/menu/:id", func(ctx *gin.Context) {
		// 获取菜单 ID
		menuID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, "无效的菜单 ID")
			return
		}

		// 调用服务层删除菜单
		_, err = c.menuService.DeleteMenuByID(menuID)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "删除菜单失败: "+err.Error())
			return
		}

		// 返回成功响应
		utils.Success(ctx, "菜单删除成功")
	})

	// 获取完整的菜单树
	menuGroup.GET("/menu/tree", middleware.RequirePermission("system:menu:tree"), func(ctx *gin.Context) {
		// 1. 获取菜单树
		menuTree, err := c.menuService.GetMenuTree()
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "获取菜单树失败: "+err.Error())
			return
		}

		// 2. 返回菜单树
		utils.Success(ctx, menuTree)
	})

	// 获取根据角色得到的菜单树
	menuGroup.GET("/menus/:roleId/tree", middleware.RequirePermission("system:menu:role:tree"), func(ctx *gin.Context) {
		// 从 URL 中获取角色 ID
		roleId, err := strconv.Atoi(ctx.Param("roleId"))
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
			return
		}

		// 1. 获取根据角色 ID 得到的菜单树
		menuTree, err := c.menuService.GetRoleMenusTreeByRoleId(roleId)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "获取菜单树失败: "+err.Error())
			return
		}

		// 2. 返回菜单树
		utils.Success(ctx, menuTree)
	})
}
