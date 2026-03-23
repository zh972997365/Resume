package routers

import (
	"log"

	"Resume/backend/config"
	"Resume/backend/handlers"
	"Resume/backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	storageService, err := services.NewStorageService(cfg.StoragePath, cfg.BaseURL, cfg.MaxFileSize)
	if err != nil {
		log.Fatalf("初始化存储服务失败: %v", err)
	}

	interviewService := services.NewInterviewService()
	companyPositionService := services.NewCompanyPositionService()
	recruitmentSourceService := services.NewRecruitmentSourceService()
	employeeService := services.NewEmployeeService()

	uploadHandler := handlers.NewUploadHandler(storageService)
	interviewHandler := handlers.NewInterviewHandler(interviewService, storageService)
	companyPositionHandler := handlers.NewCompanyPositionHandler(companyPositionService)
	recruitmentSourceHandler := handlers.NewRecruitmentSourceHandler(recruitmentSourceService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	api := r.Group("/api/v1")
	{
		api.POST("/upload", uploadHandler.UploadSingle)
		api.POST("/upload/batch", uploadHandler.UploadBatch)
		api.DELETE("/files/batch", uploadHandler.DeleteBatch)
		api.GET("/files", uploadHandler.ListFiles)
		api.GET("/files/:id", uploadHandler.GetFile)
		api.DELETE("/files/:id", uploadHandler.DeleteFile)
		api.GET("/files/search", uploadHandler.SearchFiles)
		api.GET("/files/:id/download", uploadHandler.DownloadFile)
		api.GET("/stats", uploadHandler.GetStats)
		api.POST("/rescan", uploadHandler.RescanFiles)

		interviewRoutes := api.Group("/interviews")
		{
			interviewRoutes.POST("", interviewHandler.CreateInterview)
			interviewRoutes.GET("", interviewHandler.ListInterviews)
			interviewRoutes.GET("/:id", interviewHandler.GetInterview)
			interviewRoutes.PUT("/:id", interviewHandler.UpdateInterview)
			interviewRoutes.DELETE("/:id", interviewHandler.DeleteInterview)
			interviewRoutes.POST("/upload-resume", interviewHandler.UploadResumeForInterview)
		}

		companyPositionRoutes := api.Group("/company-positions")
		{
			companyPositionRoutes.POST("", companyPositionHandler.CreatePosition)
			companyPositionRoutes.GET("", companyPositionHandler.ListPositions)
			companyPositionRoutes.GET("/:id", companyPositionHandler.GetPosition)
			companyPositionRoutes.PUT("/:id", companyPositionHandler.UpdatePosition)
			companyPositionRoutes.DELETE("/:id", companyPositionHandler.DeletePosition)
		}

		recruitmentSourceRoutes := api.Group("/recruitment-sources")
		{
			recruitmentSourceRoutes.POST("", recruitmentSourceHandler.CreateSource)
			recruitmentSourceRoutes.GET("", recruitmentSourceHandler.ListSources)
			recruitmentSourceRoutes.GET("/:id", recruitmentSourceHandler.GetSource)
			recruitmentSourceRoutes.PUT("/:id", recruitmentSourceHandler.UpdateSource)
			recruitmentSourceRoutes.DELETE("/:id", recruitmentSourceHandler.DeleteSource)
		}

		employeeRoutes := api.Group("/employees")
		{
			employeeRoutes.POST("", employeeHandler.CreateEmployee)
			employeeRoutes.GET("", employeeHandler.ListEmployees)
			employeeRoutes.GET("/:id", employeeHandler.GetEmployee)
			employeeRoutes.PUT("/:id", employeeHandler.UpdateEmployee)
			employeeRoutes.DELETE("/:id", employeeHandler.DeleteEmployee)
		}
	}

	r.Static("/uploads", cfg.StoragePath)

	r.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/interview-form.html" {
			c.HTML(200, "interview-form.html", nil)
			return
		}
		if c.Request.URL.Path == "/company-positions.html" {
			c.HTML(200, "company-positions.html", nil)
			return
		}
		if c.Request.URL.Path == "/recruitment-sources.html" {
			c.HTML(200, "recruitment-sources.html", nil)
			return
		}
		if c.Request.URL.Path == "/interviewers.html" {
			c.HTML(200, "interviewers.html", nil)
			return
		}

		c.JSON(404, gin.H{
			"success": false,
			"message": "API端点未找到或页面不存在",
		})
	})
}
