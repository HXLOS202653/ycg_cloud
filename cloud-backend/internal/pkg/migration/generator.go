// Package migration 提供数据库迁移工具的核心功能
package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Generator 迁移文件生成器
type Generator struct {
	migrationsDir string
	templateDir   string
}

// NewGenerator 创建迁移文件生成器
func NewGenerator(migrationsDir, templateDir string) *Generator {
	return &Generator{
		migrationsDir: migrationsDir,
		templateDir:   templateDir,
	}
}

// Template 迁移模板数据
type Template struct {
	Version     string
	Name        string
	ClassName   string
	TableName   string
	CreatedAt   string
	Description string
	Author      string
}

// GenerateMySQLMigration 生成MySQL迁移文件
func (g *Generator) GenerateMySQLMigration(name, description string) error {
	version := generateVersion()

	data := &Template{
		Version:     version,
		Name:        name,
		ClassName:   toPascalCase(name),
		TableName:   toSnakeCase(name),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Description: description,
		Author:      getCurrentUser(),
	}

	// 生成up文件
	upFile := filepath.Join(g.migrationsDir, "mysql", fmt.Sprintf("%s_%s.up.sql", version, name))
	if err := g.generateFromTemplate("mysql_up.tmpl", upFile, data); err != nil {
		return fmt.Errorf("生成up文件失败: %v", err)
	}

	// 生成down文件
	downFile := filepath.Join(g.migrationsDir, "mysql", fmt.Sprintf("%s_%s.down.sql", version, name))
	if err := g.generateFromTemplate("mysql_down.tmpl", downFile, data); err != nil {
		// 如果down文件生成失败，删除已生成的up文件
		_ = os.Remove(upFile)
		return fmt.Errorf("生成down文件失败: %v", err)
	}

	fmt.Printf("✅ 成功生成MySQL迁移文件:\n")
	fmt.Printf("  📄 Up:   %s\n", upFile)
	fmt.Printf("  📄 Down: %s\n", downFile)

	return nil
}

// GenerateMongoDBMigration 生成MongoDB迁移文件
func (g *Generator) GenerateMongoDBMigration(name, description string) error {
	version := generateVersion()

	data := &Template{
		Version:     version,
		Name:        name,
		ClassName:   toPascalCase(name),
		TableName:   toSnakeCase(name),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Description: description,
		Author:      getCurrentUser(),
	}

	// 生成JavaScript迁移文件
	jsFile := filepath.Join(g.migrationsDir, "mongodb", fmt.Sprintf("%s_%s.js", version, name))
	if err := g.generateFromTemplate("mongodb.tmpl", jsFile, data); err != nil {
		return fmt.Errorf("生成MongoDB迁移文件失败: %v", err)
	}

	fmt.Printf("✅ 成功生成MongoDB迁移文件:\n")
	fmt.Printf("  📄 文件: %s\n", jsFile)

	return nil
}

// generateFromTemplate 从模板生成文件
func (g *Generator) generateFromTemplate(templateName, outputFile string, data interface{}) error {
	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputFile), 0750); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 读取模板文件
	templateContent, err := g.getTemplateContent(templateName)
	if err != nil {
		return fmt.Errorf("读取模板失败: %v", err)
	}

	// 解析模板
	tmpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 创建输出文件
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer func() { _ = file.Close() }()

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	return nil
}

// getTemplateContent 获取模板内容
func (g *Generator) getTemplateContent(templateName string) (string, error) {
	// 首先尝试从文件系统读取
	templatePath := filepath.Join(g.templateDir, templateName)
	if content, err := os.ReadFile(templatePath); err == nil {
		return string(content), nil
	}

	// 如果文件不存在，使用内置模板
	switch templateName {
	case "mysql_up.tmpl":
		return mysqlUpTemplate, nil
	case "mysql_down.tmpl":
		return mysqlDownTemplate, nil
	case "mongodb.tmpl":
		return mongodbTemplate, nil
	default:
		return "", fmt.Errorf("未知的模板: %s", templateName)
	}
}

// 内置模板

const mysqlUpTemplate = `-- +migrate Up
-- 创建迁移: {{.Name}}
-- 版本: {{.Version}}
-- 描述: {{.Description}}
-- 作者: {{.Author}}
-- 创建时间: {{.CreatedAt}}

-- ============================================================================
-- 向上迁移 SQL 语句
-- ============================================================================

-- 示例: 创建表
-- CREATE TABLE {{.TableName}} (
--     id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
--     name VARCHAR(255) NOT NULL COMMENT '名称',
--     status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
--     deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
--     
--     INDEX idx_name (name),
--     INDEX idx_status (status),
--     INDEX idx_created_at (created_at),
--     INDEX idx_deleted_at (deleted_at)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='{{.Description}}';

-- 请在此处添加你的 SQL 语句
-- 注意事项:
-- 1. 使用 utf8mb4 字符集支持 emoji
-- 2. 添加适当的索引提升性能
-- 3. 使用 COMMENT 添加字段说明
-- 4. 遵循项目命名规范

`

const mysqlDownTemplate = `-- +migrate Down
-- 回滚迁移: {{.Name}}
-- 版本: {{.Version}}
-- 描述: {{.Description}}
-- 作者: {{.Author}}
-- 创建时间: {{.CreatedAt}}

-- ============================================================================
-- 向下迁移 SQL 语句 (回滚)
-- ============================================================================

-- 示例: 删除表
-- DROP TABLE IF EXISTS {{.TableName}};

-- 请在此处添加回滚 SQL 语句
-- 注意事项:
-- 1. 确保回滚操作是安全的
-- 2. 避免删除包含重要数据的表或列
-- 3. 考虑数据迁移和备份
-- 4. 测试回滚脚本的正确性

`

const mongodbTemplate = `// MongoDB 迁移: {{.Name}}
// 版本: {{.Version}}
// 描述: {{.Description}}
// 作者: {{.Author}}
// 创建时间: {{.CreatedAt}}

// ============================================================================
// MongoDB 迁移脚本
// ============================================================================

// 连接数据库
// use clouddb;

// 示例: 创建集合和索引
// db.createCollection("{{.TableName}}", {
//     validator: {
//         $jsonSchema: {
//             bsonType: "object",
//             required: ["name", "status", "created_at"],
//             properties: {
//                 name: {
//                     bsonType: "string",
//                     description: "名称字段，必填"
//                 },
//                 status: {
//                     bsonType: "string",
//                     enum: ["active", "inactive"],
//                     description: "状态字段"
//                 },
//                 created_at: {
//                     bsonType: "date",
//                     description: "创建时间"
//                 },
//                 updated_at: {
//                     bsonType: "date",
//                     description: "更新时间"
//                 }
//             }
//         }
//     }
// });

// 创建索引
// db.{{.TableName}}.createIndex({ "name": 1 });
// db.{{.TableName}}.createIndex({ "status": 1 });
// db.{{.TableName}}.createIndex({ "created_at": -1 });

// 插入初始数据
// db.{{.TableName}}.insertMany([
//     {
//         name: "示例数据",
//         status: "active",
//         created_at: new Date(),
//         updated_at: new Date()
//     }
// ]);

// 请在此处添加你的 MongoDB 操作
// 注意事项:
// 1. 使用 JSON Schema 验证器确保数据完整性
// 2. 创建适当的索引提升查询性能
// 3. 考虑分片策略 (如果使用分片)
// 4. 注意文档大小限制 (16MB)

print("✅ 迁移 {{.Name}} 执行完成");
`

// 工具函数

// generateVersion 生成版本号
func generateVersion() string {
	return time.Now().Format("20060102150405")
}

// toPascalCase 转换为帕斯卡命名
func toPascalCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if word != "" {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

// toSnakeCase 转换为蛇形命名
func toSnakeCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}

// getCurrentUser 获取当前用户
func getCurrentUser() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "system"
}

// CreateTemplateDir 创建模板目录和默认模板文件
func (g *Generator) CreateTemplateDir() error {
	templateDir := g.templateDir

	// 创建模板目录
	if err := os.MkdirAll(templateDir, 0750); err != nil {
		return fmt.Errorf("创建模板目录失败: %v", err)
	}

	// 创建默认模板文件
	templates := map[string]string{
		"mysql_up.tmpl":   mysqlUpTemplate,
		"mysql_down.tmpl": mysqlDownTemplate,
		"mongodb.tmpl":    mongodbTemplate,
	}

	for filename, content := range templates {
		filePath := filepath.Join(templateDir, filename)

		// 检查文件是否已存在
		if _, err := os.Stat(filePath); err == nil {
			continue // 文件已存在，跳过
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("创建模板文件 %s 失败: %v", filename, err)
		}

		fmt.Printf("✅ 创建模板文件: %s\n", filePath)
	}

	return nil
}
