package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"ycg_cloud/internal/utils"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
)

func main() {
	fmt.Println("=== 数据库连接测试 ===")

	// 1. 初始化配置
	fmt.Println("\n1. 初始化配置...")
	err := utils.InitConfig("", "")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	fmt.Println("✓ 配置初始化成功")

	// 2. 测试MySQL连接
	fmt.Println("\n2. 测试MySQL连接...")
	testMySQLConnection()

	// 3. 测试Redis连接
	fmt.Println("\n3. 测试Redis连接...")
	testRedisConnection()

	fmt.Println("\n=== 数据库连接测试完成 ===")
}

func testMySQLConnection() {
	// 获取数据库DSN
	dsn := utils.GetDSN()
	fmt.Printf("   数据库DSN: %s\n", dsn)

	// 尝试连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("❌ MySQL连接失败: %v\n", err)
		return
	}
	defer db.Close()

	// 设置连接池参数
	db.SetMaxOpenConns(utils.GetConfigInt("database.max_open_conns"))
	db.SetMaxIdleConns(utils.GetConfigInt("database.max_idle_conns"))
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		fmt.Printf("❌ MySQL连接测试失败: %v\n", err)
		return
	}

	fmt.Println("✓ MySQL连接成功")

	// 测试查询
	var version string
	err = db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version)
	if err != nil {
		fmt.Printf("❌ MySQL查询测试失败: %v\n", err)
		return
	}

	fmt.Printf("✓ MySQL版本: %s\n", version)
}

func testRedisConnection() {
	// 获取Redis配置
	redisAddr := utils.GetRedisAddr()
	redisPassword := utils.GetConfigString("redis.password")
	redisDB := utils.GetConfigInt("redis.database")

	fmt.Printf("   Redis地址: %s\n", redisAddr)
	fmt.Printf("   Redis数据库: %d\n", redisDB)

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	defer rdb.Close()

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("❌ Redis连接失败: %v\n", err)
		return
	}

	fmt.Printf("✓ Redis连接成功: %s\n", pong)

	// 测试基本操作
	testKey := "test_connection_key"
	testValue := "test_value"

	// 设置值
	err = rdb.Set(ctx, testKey, testValue, time.Minute).Err()
	if err != nil {
		fmt.Printf("❌ Redis设置值失败: %v\n", err)
		return
	}

	// 获取值
	val, err := rdb.Get(ctx, testKey).Result()
	if err != nil {
		fmt.Printf("❌ Redis获取值失败: %v\n", err)
		return
	}

	if val != testValue {
		fmt.Printf("❌ Redis值不匹配: 期望 %s, 实际 %s\n", testValue, val)
		return
	}

	// 删除测试键
	err = rdb.Del(ctx, testKey).Err()
	if err != nil {
		fmt.Printf("❌ Redis删除键失败: %v\n", err)
		return
	}

	fmt.Println("✓ Redis基本操作测试成功")
}
