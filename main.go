package main

import (
	"godemo/config"
	"godemo/database"
	"godemo/docs"
	"godemo/handler"
	"godemo/middleware"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Godemo API
// @version 1.0
// @description 基于Gin的用户角色权限管理系统
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}
func main() {
	// 加载配置
	config.Load()
	log.Println("✅ 配置加载成功")

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 初始化数据库
	database.InitDB()

	r := gin.Default()

	// CORS 中间件（必须在路由之前）
	r.Use(middleware.CORSMiddleware())

	docs.SwaggerInfo.BasePath = "/"

	r.POST("/api/login", handler.Login)

	// API路由组
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// 用户管理
		users := api.Group("/users")
		{
			users.GET("", middleware.PermissionMiddleware("user:view"), handler.ListUsers)
			users.GET("/:id", middleware.PermissionMiddleware("user:view"), handler.GetUser)
			users.POST("", middleware.PermissionMiddleware("user:create"), handler.CreateUser)
			users.PUT("/:id", middleware.PermissionMiddleware("user:edit"), handler.UpdateUser)
			users.DELETE("/:id", middleware.PermissionMiddleware("user:delete"), handler.DeleteUser)
		}

		// 角色管理
		roles := api.Group("/roles")
		{
			roles.GET("", middleware.PermissionMiddleware("role:view"), handler.ListRoles)
			roles.GET("/:id", middleware.PermissionMiddleware("role:view"), handler.GetRole)
			roles.POST("", middleware.PermissionMiddleware("role:create"), handler.CreateRole)
			roles.PUT("/:id", middleware.PermissionMiddleware("role:edit"), handler.UpdateRole)
			roles.DELETE("/:id", middleware.PermissionMiddleware("role:delete"), handler.DeleteRole)
		}

		// 权限管理
		api.GET("/permissions", handler.ListPermissions)
	}

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("🚀 服务启动成功，监听端口: %s", config.AppConfig.Server.Port)
	log.Printf("📚 Swagger 文档: http://localhost:%s/swagger/index.html", config.AppConfig.Server.Port)

	r.Run(":" + config.AppConfig.Server.Port)
}
