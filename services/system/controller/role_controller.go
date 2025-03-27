package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/system/dto"
	"sky_ISService/services/system/service"
	"sky_ISService/utils"
	"strconv"
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
	roleGroup := r.Group("/system")

	// 添加角色
	roleGroup.POST("/role", func(ctx *gin.Context) {
		var req dto.CreateSkySystemRoleRequest
		// 解析请求 JSON
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		role, err := c.roleService.CreateRole(req)
		if err != nil {
			// 处理 "当前角色已存在" 错误
			if err.Error() == "当前角色已存在" || err.Error() == "无法创建当前角色" {
				utils.Error(ctx, http.StatusBadRequest, err.Error()) // 400 用户请求错误
			} else {
				utils.Error(ctx, http.StatusInternalServerError, "创建用户失败: "+err.Error()) // 500 服务器错误
			}
			return
		}
		// 返回成功响应
		utils.Success(ctx, role)
	})

	// 查询单个角色
	roleGroup.GET("/role/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		role, err := c.roleService.GetRoleByID(id)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		utils.Success(ctx, role)
	})

	// 获取全部角色
	roleGroup.GET("/role", func(ctx *gin.Context) {
		// 获取请求中的分页参数
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

		// 调用封装的函数获取动态查询条件
		page, size, conditions := utils.ExtractConditions(ctx)

		// 调用 service 层获取分页数据，并传递关键字
		pagination, err := c.roleService.GetRolesWithPagination(page, size, conditions)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		// 返回分页数据
		utils.ResponseWithPagination(ctx, pagination)
	})

	// 修改角色
	roleGroup.PUT("/role", func(ctx *gin.Context) {
		var req dto.UpdateSkySystemRoleRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		updateRole, err := c.roleService.UpdateRole(req)
		if err != nil {
			// 处理 "无法修改顶级管理员账号" 错误
			if err.Error() == "无法修改顶级角色" {
				utils.Error(ctx, http.StatusBadRequest, err.Error()) // 400 用户请求错误
			} else {
				utils.Error(ctx, http.StatusInternalServerError, "更新管理员失败: "+err.Error()) // 500 服务器错误
			}
			return
		}
		// 返回更新后的管理员信息
		utils.Success(ctx, updateRole)
	})

	// 删除角色
	roleGroup.DELETE("/role/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
			return
		}
		role, err := c.roleService.DeleteRoleByID(id)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "删除角色失败:"+err.Error())
			return
		}
		utils.Success(ctx, role)
	})

	// 给角色分配可以打开菜单
	// TODO 或者查看某些菜单中的部分数据还有可以读或写的权限
	roleGroup.POST("/role/:id/menus", func(ctx *gin.Context) {
		roleID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
			return
		}
		var req dto.AssignPermissionsRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		// 检查 menuIDs 是否为空
		if len(req.MenuIDs) == 0 {
			utils.Error(ctx, http.StatusBadRequest, "请求错误: 菜单ID列表不能为空")
			return
		}
		_, err = c.roleService.AssignMenusToRole(int(roleID), req.MenuIDs)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "分配权限失败: "+err.Error())
			return
		}
		utils.Success(ctx, nil)
	})

	// TODO 给角色批量分配用户
	//roleGroup.POST("/role/:id/admins", func(ctx *gin.Context) {
	//	roleID, err := strconv.Atoi(ctx.Param("id"))
	//	if err != nil {
	//		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
	//		return
	//	}
	//
	//	var req dto.AssignAdminssRequest
	//	if err := ctx.ShouldBindJSON(&req); err != nil {
	//		utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
	//		return
	//	}
	//
	//	err = c.roleService.AssignUsersToRole(roleID, req.AdminIDs)
	//	if err != nil {
	//		utils.Error(ctx, http.StatusInternalServerError, "批量分配用户失败: "+err.Error())
	//		return
	//	}
	//
	//	utils.Success(ctx, "用户分配成功")
	//})

}
