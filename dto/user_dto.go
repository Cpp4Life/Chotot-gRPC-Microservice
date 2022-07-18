package dto

type UserUpdateDTO struct {
	Id       int    `json:"id"`
	Username string `json:"username" form:"username"`
	Address  string `json:"address" form:"address"`
	Email    string `json:"email" form:"email"`
}
