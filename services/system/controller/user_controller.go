package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/services/system/dto"
	"sky_ISService/services/system/service"
	"sky_ISService/utils"
	"strconv"
)

type AdminsController struct {
	adminsService *service.AdminsService
}

// NewUserController 构造函数，依赖注入 adminsService
func NewUserController(adminsService *service.AdminsService) *AdminsController {
	controller := &AdminsController{adminsService: adminsService}
	return controller
}

func (c *AdminsController) UserControllerRoutes(r *gin.Engine) {
	// 创建前缀的路由组
	adminGroup := r.Group("/system")

	// 添加管理员
	adminGroup.POST("/user", func(ctx *gin.Context) {
		var req dto.CreateAdminsRequest
		// 解析请求 JSON
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		// 调用服务层创建管理员
		admin, err := c.adminsService.CreateAdmin(req)
		if err != nil {
			// 处理 "当前用户已存在" 错误
			if err.Error() == "当前用户已存在" || err.Error() == "无法创建顶级管理员账号" {
				utils.Error(ctx, http.StatusBadRequest, err.Error()) // 400 用户请求错误
			} else {
				utils.Error(ctx, http.StatusInternalServerError, "创建用户失败: "+err.Error()) // 500 服务器错误
			}
			return
		}

		// 返回成功响应
		utils.Success(ctx, admin)
	})

	// 查询单个管理员用户
	adminGroup.GET("/user/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		admin, err := c.adminsService.GetAdminsByID(id)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		utils.Success(ctx, admin)
	})

	// 获取全部管理员
	adminGroup.GET("/user", func(ctx *gin.Context) {
		// 获取请求中的分页参数
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

		// 调用封装的函数获取动态查询条件
		page, size, conditions := utils.ExtractConditions(ctx)

		// 打印条件，查看提取的结果
		fmt.Println("conditions:", conditions)

		// 调用 service 层获取分页数据，并传递关键字
		pagination, err := c.adminsService.GetUsersWithPagination(page, size, conditions)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		// 返回分页数据
		utils.ResponseWithPagination(ctx, pagination)
	})

	// 修改管理员
	adminGroup.PUT("/user", func(ctx *gin.Context) {
		var req dto.UpdateAdminsRequest
		fmt.Println("req", req)
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		// 调用服务层更新管理员信息
		updatedAdmin, err := c.adminsService.UpdateAdmin(req)
		if err != nil {
			// 处理 "无法修改顶级管理员账号" 错误
			if err.Error() == "无法修改顶级管理员账号" {
				utils.Error(ctx, http.StatusBadRequest, err.Error()) // 400 用户请求错误
			} else {
				utils.Error(ctx, http.StatusInternalServerError, "更新管理员失败: "+err.Error()) // 500 服务器错误
			}
			return
		}
		// 返回更新后的管理员信息
		utils.Success(ctx, updatedAdmin)
	})

	// 删除管理员
	adminGroup.DELETE("/user/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, "无效的管理员ID")
			return
		}
		admin, err := c.adminsService.DeleteAdminByID(id)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "删除管理员失败: "+err.Error())
			return
		}
		utils.Success(ctx, admin)
	})

	// 绑定角色
	adminGroup.POST("/user/:id/roles", func(ctx *gin.Context) {
		adminID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, "无效的管理员ID")
			return
		}
		var req dto.BindRolesRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		err = c.adminsService.BindRoles(int(adminID), req.RoleIDs)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "绑定角色失败: "+err.Error())
			return
		}

		utils.Success(ctx, "角色绑定成功")
	})

	// 解绑角色
	adminGroup.DELETE("/user/:id/roles", func(ctx *gin.Context) {
		adminID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			utils.Error(ctx, http.StatusBadRequest, "无效的管理员ID")
			return
		}
		var req dto.BindRolesRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			utils.Error(ctx, http.StatusBadRequest, "请求数据错误: "+err.Error())
			return
		}
		err = c.adminsService.UnbindRoles(int(adminID), req.RoleIDs)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "解绑角色失败: "+err.Error())
			return
		}

		utils.Success(ctx, "角色解绑成功")
	})
}
