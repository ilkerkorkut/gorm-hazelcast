package main

import (
	"fmt"
	hzgorm "github.com/ilkerkorkut/gorm-hazelcast"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func insertQueryExample() {

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
	db.Save(&User{
		Username: "ilker",
		Orders:   orders,
	})
}
