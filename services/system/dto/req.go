package dto

// CreateAdminsRequest 用于创建用户
type CreateAdminsRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Notes     string `json:"notes"`
	Phone     string `json:"phone"`
	Status    bool   `json:"status"`
	CreatedBy int    `json:"created_by"`
}

// UpdateAdminsRequest 用于更新用户
type UpdateAdminsRequest struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    bool   `json:"status"`
	UpdatedBy int    `json:"updated_by"`
}
