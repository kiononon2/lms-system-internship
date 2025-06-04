//package router
//
//import (
//	"lms-system-internship/handler"
//	"lms-system-internship/repo"
//	"lms-system-internship/service"
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//	"gorm.io/gorm"
//)
//
//func SetupRoutes(db *gorm.DB, r *gin.Engine) {
//	// Repositories
//	repository := repo.NewRepository(db)
//
//	// Services
//	svc := service.NewService(repository)
//
//	// Handlers
//	courseH := handler.NewCourseHandler(svc.CourseService)
//	chapterH := handler.NewChapterHandler(svc.ChapterService)
//	lessonH := handler.NewLessonHandler(svc.LessonService)
//
//	api := r.Group("/api")
//	{
//		// Health-check
//		api.GET("", func(c *gin.Context) {
//			c.JSON(http.StatusOK, gin.H{"status": "ok"})
//		})
//
//		// Courses
//		courses := api.Group("/courses")
//		{
//			courses.GET("", courseH.GetAllCourses)
//			courses.POST("", courseH.CreateCourse)
//
//			courses.GET("/:course_id", courseH.GetCourse)
//			courses.PUT("/:course_id", courseH.UpdateCourse)
//			courses.DELETE("/:course_id", courseH.DeleteCourse)
//		}
//
//		// Chapters
//		chapters := api.Group("/chapters")
//		{
//			// Accepts `?course_id=...` as query param
//			chapters.POST("", chapterH.CreateChapter)
//			chapters.GET("", chapterH.GetAllChapters)
//			chapters.GET("/:chapter_id", chapterH.GetChapter)
//			chapters.PUT("/:chapter_id", chapterH.UpdateChapterOrder)
//			chapters.DELETE("/:chapter_id", chapterH.DeleteChapter)
//		}
//
//		// Lessons
//		lessons := api.Group("/lessons")
//		{
//			// Accepts `?chapter_id=...` as query param
//			lessons.POST("", lessonH.CreateLesson)
//			lessons.GET("", lessonH.GetAllLessons)
//			lessons.GET("/:lesson_id", lessonH.GetLesson)
//			lessons.PUT("/:lesson_id", lessonH.UpdateLessonContent)
//			lessons.DELETE("/:lesson_id", lessonH.DeleteLesson)
//		}
//
//		// Reorder lessons under a chapter
//		api.PUT("/chapters/:chapter_id/lessons/reorder", lessonH.ReorderLessons)
//	}
//}

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

		// Защищённая группа (требует JWT)
		protected := api.Group("")
		protected.Use(middleware.TokenAuthMiddleware(jwks))

		// Courses
		courses := protected.Group("/courses")
		{
			courses.GET("", courseH.GetAllCourses)

			courses.POST("", middleware.RequireRoles("ROLE_ADMIN"), courseH.CreateCourse)
			courses.GET("/:course_id", courseH.GetCourse)
			courses.PUT("/:course_id", middleware.RequireRoles("ROLE_ADMIN", "ROLE_TEACHER"), courseH.UpdateCourse)
			courses.DELETE("/:course_id", middleware.RequireRoles("ROLE_ADMIN"), courseH.DeleteCourse)
		}

		// Chapters
		chapters := protected.Group("/chapters")
		{
			chapters.POST("", middleware.RequireRoles("ROLE_ADMIN", "ROLE_TEACHER"), chapterH.CreateChapter)
			chapters.GET("", chapterH.GetAllChapters)
			chapters.GET("/:chapter_id", chapterH.GetChapter)
			chapters.PUT("/:chapter_id", middleware.RequireRoles("ROLE_ADMIN", "ROLE_TEACHER"), chapterH.UpdateChapterOrder)
			chapters.DELETE("/:chapter_id", middleware.RequireRoles("ROLE_ADMIN"), chapterH.DeleteChapter)
		}

		// Lessons
		lessons := protected.Group("/lessons")
		{
			lessons.POST("", middleware.RequireRoles("ROLE_ADMIN", "ROLE_TEACHER"), lessonH.CreateLesson)
			lessons.GET("", lessonH.GetAllLessons)
			lessons.GET("/:lesson_id", lessonH.GetLesson)
			lessons.PUT("/:lesson_id", middleware.RequireRoles("ROLE_ADMIN", "ROLE_TEACHER"), lessonH.UpdateLessonContent)
			lessons.DELETE("/:lesson_id", middleware.RequireRoles("ROLE_ADMIN"), lessonH.DeleteLesson)
		}

		protected.PUT("/chapters/:chapter_id/lessons/reorder", middleware.RequireRoles("ROLE_ADMIN", "ROLE_TEACHER"), lessonH.ReorderLessons)
	}
}
