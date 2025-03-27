package dto

// CreateAdminsRequest 用于创建用户
type CreateAdminsRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FullName  string `json:"full_name"`
	UserType  string `json:"user_type"`
	Email     string `json:"email"`
	Notes     string `json:"notes"`
	Phone     string `json:"phone"`
	Status    bool   `json:"status"`
	CreatedBy int    `json:"created_by"`
	RoleIDs   []int  `json:"role_ids"` // 角色 ID 列表
}

type UpdateAdminsRequest struct {
	ID        int    `json:"id" binding:"required"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
	UserType  string `json:"user_type"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    bool   `json:"status"`
	UpdatedBy int    `json:"updated_by"`
	Notes     string `json:"notes"`
	UpdatedAt string `json:"updated_at"`
	RoleIDs   []int  `json:"role_ids"` // 角色 ID 列表
}

// CreateSkySystemRoleRequest 用于创建角色
type CreateSkySystemRoleRequest struct {
	RoleName      string `json:"role_name" binding:"required"` // 角色名称，必填
	RoleKey       string `json:"role_key" binding:"required"`  // 角色权限字符串，必填
	RoleSort      int    `json:"role_sort" binding:"required"` // 显示顺序，必填
	Description   string `json:"description"`                  // 角色描述
	Status        bool   `json:"status" binding:"required"`    // 状态
	CreatedBy     int    `json:"created_by"`                   // 创建者 ID
	Notes         string `json:"notes"`                        // 备注
	AdminIDs      []int  `json:"admin_ids"`                    // 管理员 ID 列表
	PermissionIDs []int  `json:"permission_ids"`               // 菜单 ID 列表
}

// UpdateSkySystemRoleRequest 用于更新角色
type UpdateSkySystemRoleRequest struct {
	ID            int    `json:"id" binding:"required"` // 角色 ID，必填
	RoleName      string `json:"role_name"`             // 角色名称
	RoleKey       string `json:"role_key"`              // 角色权限字符串
	RoleSort      int    `json:"role_sort"`             // 显示顺序
	Description   string `json:"description"`           // 角色描述
	Status        bool   `json:"status"`                // 状态
	UpdatedBy     int    `json:"updated_by"`            // 更新者 ID，必填
	Notes         string `json:"notes"`                 // 备注
	AdminIDs      []int  `json:"admin_ids"`             // 管理员 ID 列表
	PermissionIDs []int  `json:"permission_ids"`        // 菜单 ID 列表
}

// CreateSkySystemMenuRequest 用于创建菜单
type CreateSkySystemMenuRequest struct {
	MenuName    string `json:"menu_name" binding:"required"`  // 菜单名称，必填
	MenuURL     string `json:"menu_url"`                      // 菜单链接（可以为空）
	ParentID    int    `json:"parent_id" binding:"required"`  // 父菜单 ID，必填
	MenuSort    int    `json:"menu_sort" binding:"required"`  // 菜单排序
	MenuType    int    `json:"menu_type" binding:"required"`  // 菜单类型: 1-目录，2-菜单，3-按钮，必填
	MenuIcon    string `json:"menu_icon"`                     // 菜单图标
	Description string `json:"description"`                   // 描述
	Status      bool   `json:"status" binding:"required"`     // 状态
	CreatedBy   int    `json:"created_by" binding:"required"` // 创建者 ID，必填
	Notes       string `json:"notes"`                         // 备注
}

// UpdateSkySystemMenuRequest 用于更新菜单
type UpdateSkySystemMenuRequest struct {
	ID          int    `json:"id" binding:"required"` // 菜单 ID，必填
	MenuName    string `json:"menu_name"`             // 菜单名称
	MenuURL     string `json:"menu_url"`              // 菜单链接（可以为空）
	ParentID    int    `json:"parent_id"`             // 父菜单 ID
	MenuSort    int    `json:"menu_sort"`             // 菜单排序
	MenuType    int    `json:"menu_type"`             // 菜单类型: 1-目录，2-菜单，3-按钮
	MenuIcon    string `json:"menu_icon"`             // 菜单图标
	Description string `json:"description"`           // 描述
	Status      bool   `json:"status"`                // 状态
	UpdatedBy   int    `json:"updated_by"`            // 更新者 ID，必填
	Notes       string `json:"notes"`                 // 备注
}

// BindRolesRequest 管理员绑定角色
type BindRolesRequest struct {
	UserID    int   `json:"user_id"`                     // 用户 ID
	RoleIDs   []int `json:"role_ids" binding:"required"` // 角色 ID 列表，必填
	CreatedAt int   `json:"created_at"`                  // // 创建时间
}

// AssignPermissionsRequest 用于角色权限分配的请求数据结构
type AssignPermissionsRequest struct {
	RoleID    int   `json:"role_id"`    // 角色 ID
	MenuIDs   []int `json:"menu_ids"`   // 菜单 ID 列表
	CreatedAt int   `json:"created_at"` // // 创建时间
}
