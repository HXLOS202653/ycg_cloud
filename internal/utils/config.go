// Package utils 提供应用程序的工具函数和配置管理功能
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ycg_cloud/internal/model"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// GlobalConfig 全局配置实例
var GlobalConfig *model.Config

// InitConfig 初始化配置
// configPath: 配置文件路径，如果为空则使用默认路径
// envFile: 环境变量文件路径，如果为空则使用默认.env文件
func InitConfig(configPath, envFile string) error {
	// 1. 加载环境变量文件
	if err := loadEnvFile(envFile); err != nil {
		return fmt.Errorf("加载环境变量文件失败: %w", err)
	}

	// 2. 设置配置文件路径
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// 3. 初始化viper
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 4. 设置环境变量前缀和自动绑定
	viper.SetEnvPrefix("YCG")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 5. 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 6. 绑定环境变量
	if err := bindEnvVars(); err != nil {
		return fmt.Errorf("绑定环境变量失败: %w", err)
	}

	// 7. 解析配置到结构体
	var config model.Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 8. 验证配置
	if err := validateConfig(&config); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// 9. 设置全局配置
	GlobalConfig = &config

	fmt.Printf("配置加载成功: %s\n", viper.ConfigFileUsed())
	return nil
}

// loadEnvFile 加载环境变量文件
func loadEnvFile(envFile string) error {
	if envFile == "" {
		// 根据环境自动选择.env文件
		env := os.Getenv("GO_ENV")
		if env == "" {
			env = "development"
		}
		envFile = fmt.Sprintf(".env.%s", env)

		// 如果环境特定的.env文件不存在，则使用默认.env文件
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			envFile = ".env"
		}
	}

	// 检查文件是否存在
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		fmt.Printf("环境变量文件不存在: %s，跳过加载\n", envFile)
		return nil
	}

	// 加载环境变量文件
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("加载环境变量文件 %s 失败: %w", envFile, err)
	}

	fmt.Printf("环境变量文件加载成功: %s\n", envFile)
	return nil
}

// bindSingleEnv 绑定单个环境变量
func bindSingleEnv(key, envVar, description string) error {
	if err := viper.BindEnv(key, envVar); err != nil {
		return fmt.Errorf("绑定%s环境变量失败: %w", description, err)
	}
	return nil
}

// bindEnvVars 绑定环境变量到viper
func bindEnvVars() error {
	// 定义环境变量映射
	envMappings := []struct {
		key, envVar, description string
	}{
		{"database.password", "YCG_DB_PASSWORD", "数据库密码"},
		{"database.host", "YCG_DB_HOST", "数据库主机"},
		{"database.port", "YCG_DB_PORT", "数据库端口"},
		{"database.username", "YCG_DB_USERNAME", "数据库用户名"},
		{"database.dbname", "YCG_DB_NAME", "数据库名称"},
		{"redis.password", "YCG_REDIS_PASSWORD", "Redis密码"},
		{"redis.host", "YCG_REDIS_HOST", "Redis主机"},
		{"redis.port", "YCG_REDIS_PORT", "Redis端口"},
		{"jwt.secret", "YCG_JWT_SECRET", "JWT密钥"},
		{"server.port", "YCG_SERVER_PORT", "服务器端口"},
		{"server.host", "YCG_SERVER_HOST", "服务器主机"},
		{"app.env", "YCG_APP_ENV", "应用环境"},
		{"app.debug", "YCG_APP_DEBUG", "应用调试"},
	}

	// 批量绑定环境变量
	for _, mapping := range envMappings {
		if err := bindSingleEnv(mapping.key, mapping.envVar, mapping.description); err != nil {
			return err
		}
	}

	return nil
}

// validateConfig 验证配置
func validateConfig(config *model.Config) error {
	if err := validateBasicConfig(config); err != nil {
		return err
	}
	if err := validateEnvType(config.App.Env); err != nil {
		return err
	}
	if err := validateLogLevel(config.Log.Level); err != nil {
		return err
	}
	return nil
}

// validateBasicConfig 验证基础配置项
func validateBasicConfig(config *model.Config) error {
	if config.App.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535之间")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("数据库主机地址不能为空")
	}
	if config.Database.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}
	return nil
}

// validateEnvType 验证环境类型
func validateEnvType(env string) error {
	validEnvs := []string{"development", "testing", "production"}
	for _, validEnv := range validEnvs {
		if env == validEnv {
			return nil
		}
	}
	return fmt.Errorf("无效的环境类型: %s，支持的类型: %v", env, validEnvs)
}

// validateLogLevel 验证日志级别
func validateLogLevel(level string) error {
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	for _, validLevel := range validLevels {
		if level == validLevel {
			return nil
		}
	}
	return fmt.Errorf("无效的日志级别: %s，支持的级别: %v", level, validLevels)
}

// GetConfig 获取全局配置
func GetConfig() *model.Config {
	return GlobalConfig
}

// GetConfigString 获取字符串配置值
func GetConfigString(key string) string {
	return viper.GetString(key)
}

// GetConfigInt 获取整数配置值
func GetConfigInt(key string) int {
	return viper.GetInt(key)
}

// GetConfigBool 获取布尔配置值
func GetConfigBool(key string) bool {
	return viper.GetBool(key)
}

// SetConfig 设置配置值（用于测试）
func SetConfig(key string, value interface{}) {
	viper.Set(key, value)
}

// GetDSN 获取数据库连接字符串
func GetDSN() string {
	if GlobalConfig == nil {
		return ""
	}

	db := GlobalConfig.Database
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.DBName,
		db.Charset,
		db.ParseTime,
		db.Loc,
	)
}

// GetRedisAddr 获取Redis地址
func GetRedisAddr() string {
	if GlobalConfig == nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", GlobalConfig.Redis.Host, GlobalConfig.Redis.Port)
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	if GlobalConfig == nil {
		return false
	}
	return GlobalConfig.App.Env == "production"
}

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	if GlobalConfig == nil {
		return true
	}
	return GlobalConfig.App.Env == "development"
}

// GetLogDir 获取日志目录
func GetLogDir() string {
	if GlobalConfig == nil {
		return "logs"
	}
	return filepath.Dir(GlobalConfig.Log.FilePath)
}
