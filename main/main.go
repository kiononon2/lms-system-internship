package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lms-system-internship/entities"
	"lms-system-internship/router"
	"log"
)

var db *gorm.DB

func initDB() {
	var err error
	//db, err = gorm.Open(postgres.Open("user=postgres password=qwerty dbname=goDB host=db port=5432 sslmode=disable"), &gorm.Config{}) - for docker
	db, err = gorm.Open(postgres.Open("user=ilia password=postgres dbname=postgres host=localhost port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&entities.Course{}, &entities.Chapter{}, &entities.Lesson{})

}

func main() {
	initDB()
	defer func() {
		s, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		s.Close()
	}()

	r := gin.Default()
	router.SetupRoutes(db, r)
	if err := r.Run(":3030"); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
