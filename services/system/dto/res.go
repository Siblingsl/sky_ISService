package dto

// AdminsResponse 返回给前端的数据
type AdminsResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   bool   `json:"status"`
}
