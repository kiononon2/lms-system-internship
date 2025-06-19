// @title LMS API
// @version 1.0
// @description Это API для системы управления курсами
// @host localhost:3030
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "lms-system-internship/docs" // важно: импорт без использования
	"lms-system-internship/entities"
	"lms-system-internship/middleware"
	"lms-system-internship/router"
	"log"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(postgres.Open("user=postgres password=qwerty dbname=goDB host=db port=5432 sslmode=disable"), &gorm.Config{})
	//db, err = gorm.Open(postgres.Open("user=ilia password=postgres dbname=postgres host=localhost port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&entities.Course{}, &entities.Chapter{}, &entities.Lesson{}, &entities.Attachment{}, &entities.LessonUser{})

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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.ErrorHandler())

	router.SetupRoutes(db, r)
	if err := r.Run(":3030"); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
