package router

import (
	"fmt"
	"lms-system-internship/handler"
	"lms-system-internship/middleware"
	"lms-system-internship/repo"
	"lms-system-internship/service"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, r *gin.Engine) {
	jwksURL := "http://keycloak:8080/realms/lms/protocol/openid-connect/certs"

	var jwks *keyfunc.JWKS
	var err error

	for i := 1; i <= 10; i++ {
		jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
			RefreshInterval: time.Hour,
			RefreshErrorHandler: func(err error) {
				fmt.Printf("JWKS refresh error: %v\n", err)
			},
		})
		if err == nil {
			break
		}
		fmt.Printf("[Retry %d/10] Failed to connect to Keycloak JWKS: %v\n", i, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to get JWKS from Keycloak: %v", err))
	}

	repository := repo.NewRepository(db)
	svc := service.NewService(repository)

	courseH := handler.NewCourseHandler(svc.CourseService)
	chapterH := handler.NewChapterHandler(svc.ChapterService)
	lessonH := handler.NewLessonHandler(svc.LessonService)

	api := r.Group("/api")
	{
		api.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		api.POST("/auth/login", handler.LoginHandler)
		api.POST("/auth/refresh", handler.RefreshTokenHandler)

		// Защищённая группа (требует JWT)
		protected := api.Group("")
		protected.Use(middleware.TokenAuthMiddleware(jwks))

		// Courses
		courses := protected.Group("/courses")
		{
			courses.GET("", courseH.GetAllCourses)

			courses.POST("", middleware.RequireRoles("ROLE_ADMIN"), courseH.CreateCourse)
			courses.GET("/:course_id", courseH.GetCourse)
			courses.PUT("/:course_id", middleware.RequireRoles("ROLE_ADMIN"), courseH.UpdateCourse)
			courses.DELETE("/:course_id", middleware.RequireRoles("ROLE_ADMIN"), courseH.DeleteCourse)
		}

		// Chapters
		chapters := protected.Group("/chapters")
		{
			chapters.POST("", middleware.RequireRoles("ROLE_ADMIN"), chapterH.CreateChapter)
			chapters.GET("", chapterH.GetAllChapters)
			chapters.GET("/:chapter_id", chapterH.GetChapter)
			chapters.PUT("/:chapter_id", middleware.RequireRoles("ROLE_ADMIN"), chapterH.UpdateChapterOrder)
			chapters.DELETE("/:chapter_id", middleware.RequireRoles("ROLE_ADMIN"), chapterH.DeleteChapter)
		}

		// Lessons
		lessons := protected.Group("/lessons")
		{
			lessons.POST("", middleware.RequireRoles("ROLE_ADMIN"), lessonH.CreateLesson)
			lessons.GET("", lessonH.GetAllLessons)
			lessons.GET("/:lesson_id", lessonH.GetLesson)
			lessons.PUT("/:lesson_id", middleware.RequireRoles("ROLE_ADMIN"), lessonH.UpdateLessonContent)
			lessons.DELETE("/:lesson_id", middleware.RequireRoles("ROLE_ADMIN"), lessonH.DeleteLesson)
		}

		protected.PUT("/chapters/:chapter_id/lessons/reorder", middleware.RequireRoles("ROLE_ADMIN"), lessonH.ReorderLessons)

		protected.POST("/admin/register", middleware.RequireRoles("ROLE_ADMIN"), handler.RegisterUser)
	}
}
