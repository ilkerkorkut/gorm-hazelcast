package main

import (
	"fmt"
	hzgorm "github.com/ilkerkorkut/gorm-hazelcast"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func queryWithPreload() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=password search_path=hazelcast sslmode=disable")
	if err != nil {
		fmt.Println("error while postgres connection !!!")
		return
	}

	db = db.Debug()

	hz, err := hzgorm.Register(db, &hzgorm.Options{
		CacheAfterPersist: true,
		Ttl:               120 * time.Second,
	})
	if err != nil {
		log.Println(err)
	}
	log.Printf("hzgorm instance %v", hz)

	var users []User

	if err := db.Table("users").Preload("Orders").Where("username = ?", "ilker1").Or("username = ?", "ilker2").Find(&users).Error; err != nil {
		log.Printf("err: %v", err)
	}
	log.Printf("query result : %v", users)
}
