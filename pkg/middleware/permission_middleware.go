package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sky_ISService/utils"
)

// RequirePermission 中间件，检查用户是否有指定的权限
func RequirePermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求的 Header 或者用户 token 获取当前用户权限
		userPermissions := getUserPermissions(ctx) // 假设 getUserPermissions 是从当前用户的上下文中获取权限的函数

		// 检查权限是否在用户的权限列表中
		if !hasPermission(userPermissions, permission) {
			utils.Error(ctx, http.StatusForbidden, "权限不足")
			ctx.Abort() // 终止请求
			return
		}

		// 权限通过，继续执行请求
		ctx.Next()
	}
}

// 获取用户权限列表的假设方法
func getUserPermissions(ctx *gin.Context) []string {
	// 假设从当前用户的 token 或上下文中获取权限
	// 实际情况可以从数据库或其他认证服务获取
	return []string{"system:menu:query", "system:user:create", "system:menu:tree", "system:menu:role:tree"}
}

// 检查是否有指定权限
func hasPermission(userPermissions []string, permission string) bool {
	for _, perm := range userPermissions {
		if perm == permission {
			return true
		}
	}
	return false
}
