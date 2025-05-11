package router

import (
	"lms-system-internship/handler"
	"lms-system-internship/repo"
	"lms-system-internship/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, r *gin.Engine) {
	// Repositories
	repository := repo.NewRepository(db)

	// Services
	svc := service.NewService(repository)

	// Handlers
	courseH := handler.NewCourseHandler(svc.CourseService)
	chapterH := handler.NewChapterHandler(svc.ChapterService)
	lessonH := handler.NewLessonHandler(svc.LessonService)

	api := r.Group("/api")
	{
		// Health-check
		api.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Courses
		courses := api.Group("/courses")
		{
			courses.GET("", courseH.GetAllCourses)
			courses.POST("", courseH.CreateCourse)

			courses.GET("/:course_id", courseH.GetCourse)
			courses.PUT("/:course_id", courseH.UpdateCourse)
			courses.DELETE("/:course_id", courseH.DeleteCourse)
		}

		// Chapters
		chapters := api.Group("/chapters")
		{
			// Accepts `?course_id=...` as query param
			chapters.POST("", chapterH.CreateChapter)
			chapters.GET("", chapterH.GetAllChapters)
			chapters.GET("/:chapter_id", chapterH.GetChapter)
			chapters.PUT("/:chapter_id", chapterH.UpdateChapterOrder)
			chapters.DELETE("/:chapter_id", chapterH.DeleteChapter)
		}

		// Lessons
		lessons := api.Group("/lessons")
		{
			// Accepts `?chapter_id=...` as query param
			lessons.POST("", lessonH.CreateLesson)
			lessons.GET("", lessonH.GetAllLessons)
			lessons.GET("/:lesson_id", lessonH.GetLesson)
			lessons.PUT("/:lesson_id", lessonH.UpdateLessonContent)
			lessons.DELETE("/:lesson_id", lessonH.DeleteLesson)
		}

		// Reorder lessons under a chapter
		api.PUT("/chapters/:chapter_id/lessons/reorder", lessonH.ReorderLessons)
	}
}
