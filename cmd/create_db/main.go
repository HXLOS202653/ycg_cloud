package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"ycg_cloud/internal/utils"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
)

func main() {
	fmt.Println("=== 数据库创建工具 ===")

	// 1. 初始化配置
	fmt.Println("\n1. 初始化配置...")
	err := utils.InitConfig("", "")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	fmt.Println("✓ 配置初始化成功")

	// 2. 创建数据库
	fmt.Println("\n2. 创建数据库...")
	createDatabase()

	fmt.Println("\n=== 数据库创建完成 ===")
}

func createDatabase() {
	// 获取数据库配置
	host := utils.GetConfigString("database.host")
	port := utils.GetConfigInt("database.port")
	username := utils.GetConfigString("database.username")
	password := utils.GetConfigString("database.password")
	dbname := utils.GetConfigString("database.dbname")

	// 构建不包含数据库名的DSN（用于连接MySQL服务器）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local",
		username, password, host, port)

	fmt.Printf("   连接MySQL服务器: %s:%d\n", host, port)
	fmt.Printf("   目标数据库: %s\n", dbname)

	// 连接到MySQL服务器
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("❌ 连接MySQL服务器失败: %v\n", err)
		return
	}
	defer db.Close()

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		fmt.Printf("❌ MySQL服务器连接测试失败: %v\n", err)
		return
	}

	fmt.Println("✓ MySQL服务器连接成功")

	// 检查数据库是否已存在
	var count int
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?",
		dbname).Scan(&count)
	if err != nil {
		fmt.Printf("❌ 检查数据库是否存在失败: %v\n", err)
		return
	}

	if count > 0 {
		fmt.Printf("✓ 数据库 '%s' 已存在\n", dbname)
		return
	}

	// 创建数据库
	createSQL := fmt.Sprintf(
		"CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		dbname)
	_, err = db.ExecContext(ctx, createSQL)
	if err != nil {
		fmt.Printf("❌ 创建数据库失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 数据库 '%s' 创建成功\n", dbname)

	// 验证数据库创建
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?",
		dbname).Scan(&count)
	if err != nil {
		fmt.Printf("❌ 验证数据库创建失败: %v\n", err)
		return
	}

	if count > 0 {
		fmt.Printf("✓ 数据库 '%s' 验证成功\n", dbname)
	} else {
		fmt.Printf("❌ 数据库 '%s' 验证失败\n", dbname)
	}
}
