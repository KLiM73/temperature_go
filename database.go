package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	ChatID   int64
	UserName string
}

func InitDB(config Config) *gorm.DB {
	log.Println("Initing DB...")
	db, err := gorm.Open(sqlite.Open(config.DbFile), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migrating DB...")
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Setting{})

	return db
}
