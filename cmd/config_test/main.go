package main

import (
	"fmt"
	"log"
	"os"

	"ycg_cloud/internal/model"
	"ycg_cloud/internal/utils"
)

// initializeConfig 初始化配置系统
func initializeConfig() {
	fmt.Println("\n1. 初始化配置...")
	if err := utils.InitConfig("", ""); err != nil {
		log.Fatal("配置初始化失败:", err)
	}
	fmt.Println("✓ 配置初始化成功")
}

// getAndValidateConfig 获取并验证配置
func getAndValidateConfig() *model.Config {
	fmt.Println("\n2. 获取配置信息...")
	config := utils.GetConfig()
	if config == nil {
		log.Fatal("获取配置失败")
	}
	return config
}

// displayAppConfig 显示应用配置
func displayAppConfig(config *model.Config) {
	fmt.Printf("\n3. 应用配置:\n")
	fmt.Printf("   应用名称: %s\n", config.App.Name)
	fmt.Printf("   版本: %s\n", config.App.Version)
	fmt.Printf("   环境: %s\n", config.App.Env)
	fmt.Printf("   调试模式: %t\n", config.App.Debug)
	fmt.Printf("   时区: %s\n", config.App.Timezone)
}

// displayServerConfig 显示服务器配置
func displayServerConfig(config *model.Config) {
	fmt.Printf("\n4. 服务器配置:\n")
	fmt.Printf("   监听地址: %s:%d\n", config.Server.Host, config.Server.Port)
	fmt.Printf("   运行模式: %s\n", config.Server.Mode)
	fmt.Printf("   读取超时: %v\n", config.Server.ReadTimeout)
	fmt.Printf("   写入超时: %v\n", config.Server.WriteTimeout)
}

// displayDatabaseConfig 显示数据库配置
func displayDatabaseConfig(config *model.Config) {
	fmt.Printf("\n5. 数据库配置:\n")
	fmt.Printf("   驱动: %s\n", config.Database.Driver)
	fmt.Printf("   主机: %s:%d\n", config.Database.Host, config.Database.Port)
	fmt.Printf("   用户名: %s\n", config.Database.Username)
	fmt.Printf("   数据库名: %s\n", config.Database.DBName)
	fmt.Printf("   字符集: %s\n", config.Database.Charset)
	fmt.Printf("   最大连接数: %d\n", config.Database.MaxOpenConns)
	fmt.Printf("   最大空闲连接数: %d\n", config.Database.MaxIdleConns)
}

// displayRedisConfig 显示Redis配置
func displayRedisConfig(config *model.Config) {
	fmt.Printf("\n6. Redis配置:\n")
	fmt.Printf("   地址: %s:%d\n", config.Redis.Host, config.Redis.Port)
	fmt.Printf("   数据库: %d\n", config.Redis.DB)
	fmt.Printf("   连接池大小: %d\n", config.Redis.PoolSize)
	fmt.Printf("   最小空闲连接: %d\n", config.Redis.MinIdleConns)
}

// displayJWTConfig 显示JWT配置
func displayJWTConfig(config *model.Config) {
	fmt.Printf("\n7. JWT配置:\n")
	fmt.Printf("   签发者: %s\n", config.JWT.Issuer)
	fmt.Printf("   过期时间: %v\n", config.JWT.ExpireTime)
	fmt.Printf("   刷新过期时间: %v\n", config.JWT.RefreshExpireTime)
	if config.JWT.Secret != "" {
		fmt.Printf("   密钥: [已设置]\n")
	} else {
		fmt.Printf("   密钥: [未设置]\n")
	}
}

// displayLogConfig 显示日志配置
func displayLogConfig(config *model.Config) {
	fmt.Printf("\n8. 日志配置:\n")
	fmt.Printf("   级别: %s\n", config.Log.Level)
	fmt.Printf("   格式: %s\n", config.Log.Format)
	fmt.Printf("   输出: %s\n", config.Log.Output)
	fmt.Printf("   文件路径: %s\n", config.Log.FilePath)
	fmt.Printf("   最大大小: %dMB\n", config.Log.MaxSize)
}

// testUtilityMethods 测试便捷方法
func testUtilityMethods() {
	fmt.Printf("\n9. 便捷方法测试:\n")
	fmt.Printf("   数据库DSN: %s\n", utils.GetDSN())
	fmt.Printf("   Redis地址: %s\n", utils.GetRedisAddr())
	fmt.Printf("   是否生产环境: %t\n", utils.IsProduction())
	fmt.Printf("   是否开发环境: %t\n", utils.IsDevelopment())
	fmt.Printf("   日志目录: %s\n", utils.GetLogDir())
}

// testEnvironmentVariables 测试环境变量覆盖
func testEnvironmentVariables() {
	fmt.Printf("\n10. 环境变量测试:\n")
	if envPort := os.Getenv("YCG_SERVER_PORT"); envPort != "" {
		fmt.Printf("   环境变量 YCG_SERVER_PORT: %s\n", envPort)
	} else {
		fmt.Printf("   环境变量 YCG_SERVER_PORT: [未设置]\n")
	}

	if envSecret := os.Getenv("YCG_JWT_SECRET"); envSecret != "" {
		fmt.Printf("   环境变量 YCG_JWT_SECRET: [已设置]\n")
	} else {
		fmt.Printf("   环境变量 YCG_JWT_SECRET: [未设置]\n")
	}
}

// testConfigGetters 测试配置获取方法
func testConfigGetters() {
	fmt.Printf("\n11. 配置获取方法测试:\n")
	fmt.Printf("   GetConfigString('app.name'): %s\n", utils.GetConfigString("app.name"))
	fmt.Printf("   GetConfigInt('server.port'): %d\n", utils.GetConfigInt("server.port"))
	fmt.Printf("   GetConfigBool('app.debug'): %t\n", utils.GetConfigBool("app.debug"))
}

func main() {
	fmt.Println("=== 配置管理系统测试 ===")

	// 初始化配置
	initializeConfig()

	// 获取配置
	config := getAndValidateConfig()

	// 显示各项配置
	displayAppConfig(config)
	displayServerConfig(config)
	displayDatabaseConfig(config)
	displayRedisConfig(config)
	displayJWTConfig(config)
	displayLogConfig(config)

	// 测试功能
	testUtilityMethods()
	testEnvironmentVariables()
	testConfigGetters()

	fmt.Println("\n=== 配置管理系统测试完成 ===")
}
