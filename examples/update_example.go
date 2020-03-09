package main

import (
	"fmt"
	hzgorm "github.com/ilkerkorkut/gorm-hazelcast"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func updateExample() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=password search_path=hazelcast sslmode=disable")
	if err != nil {
		fmt.Println("error while postgres connection !!!")
		return
	}

	db = db.Debug()

	db.AutoMigrate(&User{}, &Order{})

	hz, err := hzgorm.Register(db, &hzgorm.Options{
		CacheAfterPersist: true,
		Ttl:               120 * time.Second,
	})
	if err != nil {
		log.Println(err)
	}
	log.Printf("hzgorm instance %v", hz)

	orders := []Order{{
		Type: "Software",
	}}
	user := User{
		Username: "ilker",
		Orders:   orders,
	}

	if err := db.Save(&user).Error; err != nil {
		log.Println(err)
	}

	db.Model(&user).Update("username", "ilker2")
}
