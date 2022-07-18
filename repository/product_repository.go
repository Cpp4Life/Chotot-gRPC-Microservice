package repository

import (
	"Chotot-Microservice/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"
)

type ProductRepository interface {
	GetAllProducts() ([]entity.Product, error)
	UserProducts(id int) ([]entity.Product, error)
	InsertProduct(product *entity.Product) (*entity.Product, error)
	UpdateProduct(product *entity.Product) (*entity.Product, error)
	DeleteProduct(id string, userId string) error
	SearchProductsByName(name string, page string) ([]entity.Product, error)
}

type productConnection struct {
	conn *gorm.DB
}

func NewProductRepository(conn *gorm.DB) *productConnection {
	return &productConnection{conn: conn}
}

func (db *productConnection) GetAllProducts() ([]entity.Product, error) {
	var products []entity.Product
	if err := db.conn.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (db *productConnection) UserProducts(id int) ([]entity.Product, error) {
	var products []entity.Product
	if err := db.conn.Where("user_id = ?", id).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (db *productConnection) InsertProduct(product *entity.Product) (*entity.Product, error) {
	if err := db.conn.Create(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (db *productConnection) UpdateProduct(product *entity.Product) (*entity.Product, error) {
	output := &entity.Product{}
	if err := db.conn.Where("id = ? and user_id = ?", product.Id, product.UserId).First(&output).Error; err != nil {
		return nil, err
	}
	db.conn.Model(&output).Updates(&product)
	return output, nil
}

func (db *productConnection) DeleteProduct(id string, userId string) error {
	output := &entity.Product{}
	if err := db.conn.Where("id = ? and user_id = ?", id, userId).First(&output).Error; err != nil {
		return err
	}
	db.conn.Delete(&output)
	return nil
}

func (db *productConnection) SearchProductsByName(name string, page string) ([]entity.Product, error) {
	var (
		promotedProductPerPage   = 3
		unPromotedProductPerPage = 7
		products                 []entity.Product
		unPromotedProducts       []entity.Product
		promotedProducts         []entity.Product
	)
	numberOfPage, err := strconv.Atoi(page)
	if numberOfPage <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "page must be greater than 0")
	}

	if err != nil {
		return nil, err
	}

	if err = db.conn.Where("product_name LIKE ? AND state = 1", "%"+name+"%").Limit(promotedProductPerPage).Offset((numberOfPage - 1) * promotedProductPerPage).Find(&promotedProducts).Error; err != nil {
		return nil, err
	}

	if err = db.conn.Where("product_name LIKE ? AND state = 0", "%"+name+"%").Limit(unPromotedProductPerPage).Offset((numberOfPage - 1) * unPromotedProductPerPage).Find(&unPromotedProducts).Error; err != nil {
		return nil, err
	}

	products = append(products, unPromotedProducts...)
	products = append(products, promotedProducts...)
	return products, nil
}
