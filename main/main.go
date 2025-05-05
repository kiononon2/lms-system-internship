package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lms-system-internship/entities"
	"log"
	"net/http"
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

	//repository := repo.NewRepository(db)
	//
	//services := service.NewService(repository)
	//
	////course := entities.Course{
	////	ID:          0,
	////	Name:        "Example Course",
	////	Description: "example course",
	////	CreatedAt:   time.Time{},
	////	UpdatedAt:   time.Time{},
	////	Chapters:    nil,
	////}
	//
	//services.CourseService.DeleteCourse(context.Background(), 1)

	server := &http.Server{
		Addr: ":3030",
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
