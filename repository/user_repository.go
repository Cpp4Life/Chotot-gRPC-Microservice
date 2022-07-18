package repository

import (
	"Chotot-Microservice/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	VerifyCredential(phone string) (*entity.User, error)
	IsDuplicatePhone(phone string) (bool, error)
	InsertUser(user *entity.User) (*entity.User, error)
	UserProfile(id int) (*entity.User, error)
	UpdateUser(id int, user *entity.User) (*entity.User, error)
}

type userConnection struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) *userConnection {
	return &userConnection{conn: conn}
}

func (db *userConnection) VerifyCredential(phone string) (*entity.User, error) {
	user := &entity.User{}
	if err := db.conn.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (db *userConnection) IsDuplicatePhone(phone string) (bool, error) {
	user := &entity.User{}
	if err := db.conn.Where("phone = ?", phone).First(&user).Error; err != nil {
		return false, err
	}
	if user.Phone == "" {
		return false, nil
	}
	return true, nil
}

func (db *userConnection) InsertUser(user *entity.User) (*entity.User, error) {
	user.Passwd = hashAndSalt([]byte(user.Passwd))
	if err := db.conn.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (db *userConnection) UserProfile(id int) (*entity.User, error) {
	user := &entity.User{}
	if err := db.conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (db *userConnection) UpdateUser(id int, user *entity.User) (*entity.User, error) {
	output := &entity.User{}
	if err := db.conn.Where("id = ?", id).First(&output).Error; err != nil {
		return nil, err
	}
	db.conn.Model(&output).Updates(&user)
	return output, nil
}

func hashAndSalt(passwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(passwd, bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}
	return string(hash)
}
