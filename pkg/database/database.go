package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kais-blkc/go-blog/internal/model"
)

func Connect(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
