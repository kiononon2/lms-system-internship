package router

import (
	"lms-system-internship/handler"
	"lms-system-internship/service"
	"lms-system-internship/storage/repo"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, r *gin.Engine) {
	// Repositories
	repo := repo.NewRepository(db)

	// Services
	svc := service.NewService(repo)

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

			byCourse := courses.Group("/:course_id")
			{
				byCourse.GET("", courseH.GetCourse)
				byCourse.PUT("", courseH.UpdateCourse)
				byCourse.DELETE("", courseH.DeleteCourse)

				// Optional: Add chapterH.GetChaptersByCourseID if implemented
			}
		}

		// Chapters
		chapters := api.Group("/chapters")
		{
			chapters.POST("", chapterH.CreateChapter)

			byChapter := chapters.Group("/:chapter_id")
			{
				//byChapter.GET("", chapterH.GetChapter)
				byChapter.PUT("", chapterH.UpdateChapter)
				byChapter.DELETE("", chapterH.DeleteChapter)
				// Optional: GET /lessons if needed
			}
		}

		// Lessons
		lessons := api.Group("/lessons")
		{
			lessons.POST("", lessonH.CreateLesson)

			byLesson := lessons.Group("/:lesson_id")
			{
				//byLesson.GET("", lessonH.GetLesson)
				byLesson.PUT("", lessonH.UpdateLesson)
				byLesson.DELETE("", lessonH.DeleteLesson)
			}
		}
	}
}
