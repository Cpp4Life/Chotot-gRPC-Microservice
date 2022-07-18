package dto

type RegisterDTO struct {
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Passwd   string `json:"password" form:"password" binding:"required"`
	Address  string `json:"address" form:"address"`
	Email    string `json:"email" form:"email"`
	IsAdmin  bool   `json:"is_admin" form:"is_admin"`
}
