package dto

import "time"

// SkySystemAdminsResponse 返回给前端的数据
type SkySystemAdminsResponse struct {
	ID        int       `json:"id"`         // 使用 int 类型
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"password"`   // 密码
	FullName  string    `json:"full_name"`  // 全名
	UserType  string    `json:"user_type"`  // 管理员类型（00系统管理员）
	Email     string    `json:"email"`      // 邮箱
	Phone     string    `json:"phone"`      // 手机号
	Status    bool      `json:"status"`     // 状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
	CreatedBy int       `json:"created_by"` // 创建者 ID
	UpdatedBy int       `json:"updated_by"` // 修改者ID
	Notes     string    `json:"notes"`      // 备注
	RolesID   []int     `json:"roles_id"`   // 角色 ID 列表
}

// SkySystemRolesResponse 用于返回角色数据
type SkySystemRolesResponse struct {
	ID          int       `json:"id"`          // 角色 ID
	RoleName    string    `json:"role_name"`   // 角色名称
	RoleKey     string    `json:"role_key"`    // 角色权限字符串
	RoleSort    int       `json:"role_sort"`   // 显示顺序
	Description string    `json:"description"` // 角色描述
	Status      bool      `json:"status"`      // 状态
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
	CreatedBy   int       `json:"created_by"`  // 创建者 ID
	UpdatedBy   int       `json:"updated_by"`  // 修改者ID
	Notes       string    `json:"notes"`       // 备注
}

// SkySystemMenuResponse 用于返回菜单数据
type SkySystemMenuResponse struct {
	ID          int    `json:"id"`          // 菜单 ID
	MenuName    string `json:"menu_name"`   // 菜单名称
	MenuURL     string `json:"menu_url"`    // 菜单链接
	ParentID    int    `json:"parent_id"`   // 父菜单 ID
	MenuSort    int    `json:"menu_sort"`   // 菜单排序
	MenuType    int    `json:"menu_type"`   // 菜单类型: 1-目录，2-菜单，3-按钮
	MenuIcon    string `json:"menu_icon"`   // 菜单图标
	Description string `json:"description"` // 描述
	CreatedAt   string `json:"created_at"`  // 创建时间
	UpdatedAt   string `json:"updated_at"`  // 更新时间
}

// MenuItem 代表一个菜单项
type MenuItem struct {
	ID       int        `json:"id"`        // 菜单 ID
	MenuName string     `json:"menu_name"` // 菜单名称
	MenuURL  string     `json:"menu_url"`  // 菜单链接
	ParentID int        `json:"parent_id"` // 父菜单 ID
	MenuSort int        `json:"menu_sort"` // 菜单排序
	MenuType int        `json:"menu_type"` // 菜单类型: 1-目录，2-菜单，3-按钮
	MenuIcon string     `json:"menu_icon"` // 菜单图标
	Children []MenuItem `json:"children"`  // 子菜单
}
