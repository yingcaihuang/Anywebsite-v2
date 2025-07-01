package main

import (
	"html/template"
	"log"
	"static-hosting-server/internal/api"
	"static-hosting-server/internal/config"
	"static-hosting-server/internal/database"
	"static-hosting-server/internal/scheduler"
	"static-hosting-server/internal/web"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由器
	router := gin.Default()

	// 设置自定义模板函数
	router.SetFuncMap(template.FuncMap{
		"add":      func(a, b int) int { return a + b },
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	})

	// 静态文件服务
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// 设置路由
	api.SetupRoutes(router, db, cfg)
	web.SetupRoutes(router, db, cfg)

	// 启动定时任务
	scheduler.Start(db)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
