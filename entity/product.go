package entity

import "time"

type Product struct {
	Id          int       `json:"id" gorm:"primary_key;auto_increment"`
	ProductName string    `json:"product_name" gorm:"type:varchar(255);not null;unique"`
	UserId      int       `json:"user_id" gorm:"type:int"`
	CatId       string    `json:"cat_id" gorm:"type:VARCHAR(10)"`
	TypeId      string    `json:"type_id" gorm:"type:VARCHAR(10)"`
	Price       float64   `json:"price" gorm:"type:float"`
	State       bool      `json:"state" gorm:"type:boolean"`
	CreatedTime time.Time `json:"created_time" gorm:"type:timestamp"`
	ExpiredTime time.Time `json:"expired_time" gorm:"type:timestamp"`
	Address     string    `json:"address" gorm:"type:varchar(255)"`
	Content     string    `json:"content" gorm:"type:varchar(255)"`
}
