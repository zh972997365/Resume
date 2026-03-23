package main

import (
	"log"
	"os"
	"time"

	"Resume/backend/config"
	"Resume/backend/database"
	"Resume/backend/middleware"
	"Resume/backend/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if err := os.MkdirAll(cfg.StoragePath, 0755); err != nil {
		log.Fatalf("创建存储目录失败: %v", err)
	}

	log.Printf("配置文件加载成功:")
	log.Printf("  服务器地址: %s:%s", cfg.Host, cfg.Port)
	log.Printf("  存储路径: %s", cfg.StoragePath)
	log.Printf("  基础URL: %s", cfg.BaseURL)
	log.Printf("  调试模式: %v", cfg.Debug)

	database.InitDB(cfg)

	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(middleware.CORS())

	// 设置 Go HTML 模板的自定义分隔符，避免与 Vue.js 的 {{}} 冲突
	r.Delims("[[", "]]") // 新增此行

	r.Static("/static", "./frontend/static")
	r.LoadHTMLGlob("./frontend/templates/*")

	routers.SetupRoutes(r, cfg)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "简历上传系统",
		})
	})

	r.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "简历上传系统",
		})
	})

	r.GET("/interview-form.html", func(c *gin.Context) {
		c.HTML(200, "interview-form.html", gin.H{
			"title": "面试管理表单",
		})
	})

	r.GET("/company-positions.html", func(c *gin.Context) {
		c.HTML(200, "company-positions.html", gin.H{
			"title": "岗位管理",
		})
	})

	r.GET("/recruitment-sources.html", func(c *gin.Context) {
		c.HTML(200, "recruitment-sources.html", gin.H{
			"title": "招聘来源管理",
		})
	})

	r.GET("/interviewers.html", func(c *gin.Context) {
		c.HTML(200, "interviewers.html", gin.H{
			"title": "面试官管理",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":       "ok",
			"service":      "简历上传系统",
			"version":      "1.0.0",
			"storage_path": cfg.StoragePath,
			"debug":        cfg.Debug,
		})
	})

	r.GET("/api/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success":   true,
			"message":   "API服务正常运行",
			"timestamp": time.Now().Unix(),
		})
	})

	log.Printf("🚀 服务器启动在 http://%s:%s", cfg.Host, cfg.Port)
	log.Printf("📁 文件存储目录: %s", cfg.StoragePath)
	log.Printf("📊 前端页面地址: http://%s:%s/", cfg.Host, cfg.Port)

	if err := r.Run(cfg.Host + ":" + cfg.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
