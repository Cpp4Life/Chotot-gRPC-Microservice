package entity

type User struct {
	Id       int    `json:"id" gorm:"primary_key;auto_increment"`
	Phone    string `json:"phone" gorm:"type:varchar(12);not null;unique"`
	Username string `json:"username" gorm:"type:varchar(50);not null"`
	Passwd   string `json:"passwd" gorm:"type:varchar(255);not null"`
	Address  string `json:"address" gorm:"type:varchar(255)"`
	Email    string `json:"email" gorm:"type:varchar(50);unique"`
	IsAdmin  bool   `json:"is_admin" gorm:"type:boolean;default:0"`
	Token    string `json:"token" gorm:"-"`
}
