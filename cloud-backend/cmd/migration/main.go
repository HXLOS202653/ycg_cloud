// Package main provides the database migration tool for the YCG Cloud system.
// It supports MySQL and MongoDB database schema migrations with version control and rollback capabilities.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// 数据库连接配置
	dsn        = flag.String("dsn", "", "数据库连接字符串")
	configFile = flag.String("config", "configs/config.yaml", "配置文件路径")

	// 迁移相关选项
	migrationsDir = flag.String("migrations-dir", "migrations", "迁移文件目录")
	version       = flag.String("version", "", "指定迁移版本")
	steps         = flag.Int("steps", 0, "迁移步数")
	force         = flag.Bool("force", false, "强制执行迁移")
	dryRun        = flag.Bool("dry-run", false, "预览模式，不实际执行")
	verbose       = flag.Bool("verbose", false, "详细输出")
)

// MigrationTool 数据库迁移工具
type MigrationTool struct {
	db            *gorm.DB
	migrationsDir string
	verbose       bool
	dryRun        bool
}

// NewMigrationTool 创建迁移工具实例
func NewMigrationTool(dsn, migrationsDir string, verbose, dryRun bool) (*MigrationTool, error) {
	// 配置GORM日志级别
	logLevel := logger.Silent
	if verbose {
		logLevel = logger.Info
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return &MigrationTool{
		db:            db,
		migrationsDir: migrationsDir,
		verbose:       verbose,
		dryRun:        dryRun,
	}, nil
}

func main() {
	// 设置根命令
	rootCmd := &cobra.Command{
		Use:   "migration",
		Short: "网络云盘系统数据库迁移工具",
		Long: `网络云盘系统数据库迁移工具
支持MySQL和MongoDB的数据库结构迁移，包括表创建、索引管理、数据转换等功能。
遵循版本控制和回滚机制，确保数据库结构的安全迁移。`,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			// 解析标志
			flag.Parse()

			// 验证必要参数
			if *dsn == "" {
				fmt.Println("错误: 必须指定数据库连接字符串 (--dsn)")
				os.Exit(1)
			}

			// 验证迁移目录
			if _, err := os.Stat(*migrationsDir); os.IsNotExist(err) {
				fmt.Printf("错误: 迁移目录不存在: %s\n", *migrationsDir)
				os.Exit(1)
			}
		},
	}

	// 添加持久化标志
	rootCmd.PersistentFlags().StringVar(dsn, "dsn", "", "数据库连接字符串 (必需)")
	rootCmd.PersistentFlags().StringVar(configFile, "config", "configs/config.yaml", "配置文件路径")
	rootCmd.PersistentFlags().StringVar(migrationsDir, "migrations-dir", "migrations", "迁移文件目录")
	rootCmd.PersistentFlags().BoolVar(verbose, "verbose", false, "详细输出")
	rootCmd.PersistentFlags().BoolVar(dryRun, "dry-run", false, "预览模式，不实际执行")
	rootCmd.PersistentFlags().BoolVar(force, "force", false, "强制执行迁移")

	// up 命令 - 执行迁移
	upCmd := &cobra.Command{
		Use:   "up",
		Short: "执行数据库迁移",
		Long: `执行数据库迁移，将数据库结构升级到最新版本或指定版本。
支持指定版本号或步数来控制迁移范围。`,
		Example: `  # 执行所有待迁移的版本
  migration up --dsn "user:password@tcp(localhost:3306)/clouddb?charset=utf8mb4&parseTime=True&loc=Local"
  
  # 执行到指定版本
  migration up --version 003 --dsn "..."
  
  # 执行指定步数
  migration up --steps 2 --dsn "..."
  
  # 预览模式
  migration up --dry-run --dsn "..."`,
		Run: func(_ *cobra.Command, _ []string) {
			tool, err := NewMigrationTool(*dsn, *migrationsDir, *verbose, *dryRun)
			if err != nil {
				log.Fatalf("初始化迁移工具失败: %v", err)
			}

			if err := tool.Up(*version, *steps); err != nil {
				log.Fatalf("执行迁移失败: %v", err)
			}
		},
	}
	upCmd.Flags().StringVar(version, "version", "", "迁移到指定版本")
	upCmd.Flags().IntVar(steps, "steps", 0, "执行指定步数的迁移")

	// down 命令 - 回滚迁移
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "回滚数据库迁移",
		Long: `回滚数据库迁移，将数据库结构回退到之前的版本。
支持指定版本号或步数来控制回滚范围。`,
		Example: `  # 回滚一个版本
  migration down --dsn "user:password@tcp(localhost:3306)/clouddb?charset=utf8mb4&parseTime=True&loc=Local"
  
  # 回滚到指定版本
  migration down --version 001 --dsn "..."
  
  # 回滚指定步数
  migration down --steps 2 --dsn "..."`,
		Run: func(_ *cobra.Command, _ []string) {
			tool, err := NewMigrationTool(*dsn, *migrationsDir, *verbose, *dryRun)
			if err != nil {
				log.Fatalf("初始化迁移工具失败: %v", err)
			}

			if err := tool.Down(*version, *steps); err != nil {
				log.Fatalf("回滚迁移失败: %v", err)
			}
		},
	}
	downCmd.Flags().StringVar(version, "version", "", "回滚到指定版本")
	downCmd.Flags().IntVar(steps, "steps", 1, "回滚指定步数")

	// status 命令 - 查看迁移状态
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "查看数据库迁移状态",
		Long: `查看当前数据库的迁移状态，包括已执行的迁移版本和待执行的迁移。
显示每个迁移文件的执行状态和时间。`,
		Run: func(_ *cobra.Command, _ []string) {
			tool, err := NewMigrationTool(*dsn, *migrationsDir, *verbose, *dryRun)
			if err != nil {
				log.Fatalf("初始化迁移工具失败: %v", err)
			}

			if err := tool.Status(); err != nil {
				log.Fatalf("查看迁移状态失败: %v", err)
			}
		},
	}

	// create 命令 - 创建迁移文件
	createCmd := &cobra.Command{
		Use:   "create [name]",
		Short: "创建新的迁移文件",
		Long: `创建新的迁移文件，包括up和down两个SQL文件。
文件名会自动添加时间戳前缀以确保版本顺序。`,
		Args: cobra.ExactArgs(1),
		Example: `  # 创建用户表迁移
  migration create create_users_table
  
  # 创建索引迁移
  migration create add_user_indexes`,
		Run: func(_ *cobra.Command, args []string) {
			tool, err := NewMigrationTool(*dsn, *migrationsDir, *verbose, *dryRun)
			if err != nil {
				log.Fatalf("初始化迁移工具失败: %v", err)
			}

			migrationName := args[0]
			if err := tool.Create(migrationName); err != nil {
				log.Fatalf("创建迁移文件失败: %v", err)
			}
		},
	}

	// version 命令 - 显示版本信息
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "显示工具版本信息",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("网络云盘系统数据库迁移工具 v1.0.0")
			fmt.Println("支持 MySQL 8.0+ 和 MongoDB 4.4+")
			fmt.Println("Build: " + getBuildInfo())
		},
	}

	// 添加子命令
	rootCmd.AddCommand(upCmd, downCmd, statusCmd, createCmd, versionCmd)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("命令执行失败: %v", err)
	}
}

// getBuildInfo 获取构建信息
func getBuildInfo() string {
	// 这里可以在构建时通过 ldflags 注入版本信息
	return "dev-build"
}

// 工具辅助函数

// ensureDir 确保目录存在
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0750)
	}
	return nil
}

// getMigrationsPath 获取迁移文件路径
func getMigrationsPath(migrationsDir, dbType string) string {
	return filepath.Join(migrationsDir, dbType)
}

// validateMigrationName 验证迁移名称
func validateMigrationName(name string) error {
	if name == "" {
		return fmt.Errorf("迁移名称不能为空")
	}

	// 检查是否包含特殊字符
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return fmt.Errorf("迁移名称只能包含字母、数字和下划线")
		}
	}

	return nil
}
