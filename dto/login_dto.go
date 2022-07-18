package dto

type LoginDTO struct {
	Phone  string `json:"phone" form:"phone" binding:"required"`
	Passwd string `json:"password" form:"password" binding:"required"`
}
