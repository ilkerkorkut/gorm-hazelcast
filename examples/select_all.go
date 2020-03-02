package main

import (
	"fmt"
	hzgorm "github.com/ilkerkorkut/gorm-hazelcast"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func selectAllQuery() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=password search_path=hazelcast sslmode=disable")
	if err != nil {
		fmt.Println("error while postgres connection !!!")
		return
	}

	db = db.Debug()

	_, err = hzgorm.Register(db, &hzgorm.Options{
		CacheAfterPersist: true,
		Ttl:               120 * time.Second,
	})

	var users []User

	if err := db.Preload("Orders").Find(&users).Error; err != nil {
		log.Printf("err: %v", err)
	}

	log.Printf("Result: %v", users)

}
