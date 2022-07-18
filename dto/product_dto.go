package dto

type ProductDTO struct {
	Id          string  `json:"id" form:"id"`
	ProductName string  `json:"product_name" form:"product_name"`
	UserId      int     `json:"user_id" form:"user_id"`
	CatId       string  `json:"cat_id" form:"cat_id"`
	TypeId      string  `json:"type_id" form:"type_id"`
	Price       float64 `json:"price" form:"price"`
	State       bool    `json:"state" form:"state"`
	Address     string  `json:"address" form:"address"`
	Content     string  `json:"content" form:"content"`
}
