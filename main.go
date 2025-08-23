package main

import (
	"fmt"
	"log"
	"net/http"

	"ycg_cloud/internal/utils"

	"github.com/gin-gonic/gin"
)

// main 程序入口点，启动Gin HTTP服务器
func main() {
	// 初始化配置
	if err := utils.InitConfig("", ""); err != nil {
		log.Fatal("配置初始化失败:", err)
	}

	// 获取配置
	config := utils.GetConfig()
	if config == nil {
		log.Fatal("获取配置失败")
	}

	// 设置Gin模式
	gin.SetMode(config.Server.Mode)

	// 创建Gin引擎
	router := gin.Default()

	// 添加基础中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 添加CORS中间件
	router.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	})

	// 基础路由
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "服务运行正常",
			"data":    gin.H{"status": "healthy"},
		})
	})

	// 启动服务器
	log.Printf("%s 后端服务启动中...", config.App.Name)
	log.Printf("版本: %s", config.App.Version)
	log.Printf("环境: %s", config.App.Env)
	log.Printf("调试模式: %t", config.App.Debug)
	serverAddr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("服务地址: http://%s", serverAddr)
	log.Printf("健康检查: http://%s/api/v1/health", serverAddr)

	if err := router.Run(fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}
