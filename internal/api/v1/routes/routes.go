package routes

import (
	"time"

	"github.com/0xBoji/web3-edu-core/internal/api/middleware"
	"github.com/0xBoji/web3-edu-core/internal/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the v1 API routes
func RegisterRoutes(router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/api/v1")

	// Apply rate limiting to auth routes
	authRateLimiter := middleware.RateLimitMiddleware(100, time.Minute)

	// Auth routes - public endpoints
	authHandler := handlers.NewAuthHandler()
	auth := v1.Group("/auth")
	auth.Use(authRateLimiter)
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		userHandler := handlers.NewUserHandler()

		// User profile routes
		users := protected.Group("/users")
		{
			// Current user profile
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
			// TODO: Implement these endpoints
			// users.PATCH("/me/password", userHandler.UpdatePassword)
			// users.GET("/me/enrollments", userHandler.GetEnrollments)

			// Admin routes for user management
			users.GET("", middleware.RoleMiddleware("admin"), userHandler.List)
			users.GET("/:id", middleware.RoleMiddleware("admin"), userHandler.Get)
			users.PUT("/:id", middleware.RoleMiddleware("admin"), userHandler.Update)
			users.DELETE("/:id", middleware.RoleMiddleware("admin"), userHandler.Delete)
		}

		// Category routes
		categoryHandler := handlers.NewCategoryHandler()
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.List)
			categories.GET("/:id", categoryHandler.Get)
		}

		// Admin category routes
		adminCategories := protected.Group("/admin/categories")
		adminCategories.Use(middleware.RoleMiddleware("admin"))
		{
			adminCategories.POST("", categoryHandler.Create)
			adminCategories.PUT("/:id", categoryHandler.Update)
			adminCategories.DELETE("/:id", categoryHandler.Delete)
		}

		// Course routes
		// courseHandler := handlers.NewCourseHandler()
		// courses := v1.Group("/courses")
		// {
		//     courses.GET("", courseHandler.List)
		//     courses.GET("/featured", courseHandler.GetFeatured)
		//     courses.GET("/:id", courseHandler.Get)
		//     courses.GET("/:id/lessons", courseHandler.GetLessons)
		//     courses.GET("/:id/reviews", courseHandler.GetReviews)
		// }
		//
		// // Protected course routes
		// protectedCourses := protected.Group("/courses")
		// {
		//     protectedCourses.POST("/:id/enroll", courseHandler.Enroll)
		//     protectedCourses.POST("/:id/reviews", courseHandler.AddReview)
		// }
		//
		// // Admin course routes
		// adminCourses := protected.Group("/admin/courses")
		// adminCourses.Use(middleware.RoleMiddleware("admin", "instructor"))
		// {
		//     adminCourses.POST("", courseHandler.Create)
		//     adminCourses.PUT("/:id", courseHandler.Update)
		//     adminCourses.DELETE("/:id", courseHandler.Delete)
		// }

		// Lesson routes
		// lessonHandler := handlers.NewLessonHandler()
		// lessons := protected.Group("/lessons")
		// {
		//     lessons.GET("/:id", lessonHandler.Get)
		//     lessons.GET("/:id/progress", lessonHandler.GetProgress)
		//     lessons.POST("/:id/progress", lessonHandler.UpdateProgress)
		//     lessons.POST("/:id/complete", lessonHandler.Complete)
		// }
		//
		// // Admin lesson routes
		// adminLessons := protected.Group("/admin/lessons")
		// adminLessons.Use(middleware.RoleMiddleware("admin", "instructor"))
		// {
		//     adminLessons.POST("", lessonHandler.Create)
		//     adminLessons.PUT("/:id", lessonHandler.Update)
		//     adminLessons.DELETE("/:id", lessonHandler.Delete)
		// }

		// Enrollment routes
		// enrollmentHandler := handlers.NewEnrollmentHandler()
		// enrollments := protected.Group("/enrollments")
		// {
		//     enrollments.GET("", enrollmentHandler.List)
		//     enrollments.POST("", enrollmentHandler.Create)
		//     enrollments.DELETE("/:id", enrollmentHandler.Delete)
		// }

		// Progress routes
		// progressHandler := handlers.NewProgressHandler()
		// progress := protected.Group("/progress")
		// {
		//     progress.GET("", progressHandler.List)
		//     progress.POST("", progressHandler.Create)
		//     progress.PUT("/:id", progressHandler.Update)
		// }

		// I18n routes
		i18nHandler := handlers.NewI18nHandler()
		i18n := v1.Group("/i18n")
		{
			i18n.GET("/:language", i18nHandler.GetTranslations)
		}
	}
}
