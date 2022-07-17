package user

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	fmt.Println("test print new Repository:", &repository{db}, "db: ", db)
	isiDb := *db
	isiRepo := &repository{db}
	fmt.Println("isi db:", isiDb)
	fmt.Println("isi repository:", *isiRepo)
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
