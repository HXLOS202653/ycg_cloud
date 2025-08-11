# 网络云盘系统 - API设计规范

## 📋 概述

API设计规范定义了网络云盘系统中所有API接口的设计标准，基于现代化技术栈（React 18 + TypeScript + Go + Gin + MySQL + MongoDB + Redis）构建的企业级云存储解决方案。本规范涵盖URL结构、HTTP方法、请求/响应格式、错误处理、认证授权等各个方面，确保API的一致性、可维护性、安全性和高性能。

### 系统架构概览
- **前端**: React 18 + TypeScript + Ant Design + Redux Toolkit
- **后端**: Go 1.24.6 + Gin 1.10.0 + GORM 1.30.1 + WebSocket
- **数据库**: MySQL 8.0 (主数据库) + MongoDB 6.0 (文档数据) + Redis 7.0 (缓存)
- **存储**: MinIO/S3 对象存储 + MySQL 全文索引搜索
- **部署**: Docker + Kubernetes + Nginx

## 🏗️ API设计原则

### 1. 核心原则
- **RESTful设计**: 遵循REST架构风格，资源导向的URL设计
- **一致性**: 保持命名、结构、行为的一致性，统一的响应格式
- **简洁性**: API设计简单易懂，减少学习成本，直观的接口设计
- **可扩展性**: 预留扩展空间，支持版本演进，模块化设计
- **安全性**: 严格的认证授权和数据验证，防止常见攻击
- **性能优化**: 支持分页、过滤、字段选择等优化手段，缓存策略

### 2. 设计理念
```
面向资源设计 + 语义化操作 + 统一响应格式 + 完善错误处理 + 安全第一 + 性能优化
```

### 3. 字段映射转换规范

#### 3.1 命名转换规则

**转换链路**: `数据库字段(snake_case) ↔ Go结构体(PascalCase) ↔ JSON(camelCase) ↔ 前端(camelCase)`

**转换规则详细说明**:
```yaml
# 命名转换映射表
database_field: file_name         # MySQL: snake_case (下划线分隔)
go_struct_field: FileName         # Go: PascalCase (首字母大写驼峰)
json_field: fileName              # JSON/API: camelCase (小写驼峰)
frontend_field: fileName          # Frontend: camelCase (小写驼峰)

# 特殊字段处理
id: id                            # ID字段保持小写
created_at: createdAt             # 时间字段统一使用camelCase
updated_at: updatedAt
deleted_at: deletedAt
user_id: userId                   # 外键字段遵循camelCase规则
```

**Go结构体标签规范**:
```go
type File struct {
    ID          int64     `json:"id" db:"id"`
    FileName    string    `json:"fileName" db:"file_name"`
    FileSize    int64     `json:"fileSize" db:"file_size"`
    ContentType string    `json:"contentType" db:"content_type"`
    UserID      int64     `json:"userId" db:"user_id"`
    CreatedAt   time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
```

#### 3.2 自动化代码生成工具

**选定工具链**:
```yaml
# 1. 数据库 → Go结构体
tool: sqlc                        # SQL代码生成器
config: .sqlc/sqlc.yaml
command: sqlc generate

# 2. Go结构体 → OpenAPI文档  
tool: swaggo/swag                 # Swagger文档生成
command: swag init -g cmd/server/main.go

# 3. OpenAPI → TypeScript类型
tool: openapi-typescript          # TypeScript类型生成
command: npx openapi-typescript api.yaml -o types/api.ts

# 4. 字段映射一致性检查
tool: 自定义脚本 (mapping-validator)
command: go run tools/mapping-validator.go
```

**SQLC配置示例**:
```yaml
# .sqlc/sqlc.yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "./internal/repository/queries"
    schema: "./migrations"
    gen:
      go:
        package: "repository"
        out: "./internal/repository"
        sql_package: "database/sql"
        emit_json_tags: true
        emit_db_tags: true
        json_tags_case_style: "camel"        # JSON字段使用camelCase
        output_file_pattern: "{{.Table}}.sql.go"
```

**字段映射一致性检查工具**:
```go
// tools/mapping-validator.go
package main

import (
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "strings"
)

// 检查Go结构体字段映射一致性
func validateFieldMapping(structName string) error {
    // 1. 解析Go结构体
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, "models.go", nil, parser.ParseComments)
    if err != nil {
        return err
    }
    
    // 2. 提取字段标签
    for _, decl := range node.Decls {
        if genDecl, ok := decl.(*ast.GenDecl); ok {
            for _, spec := range genDecl.Specs {
                if typeSpec, ok := spec.(*ast.TypeSpec); ok {
                    if typeSpec.Name.Name == structName {
                        return validateStruct(typeSpec)
                    }
                }
            }
        }
    }
    
    return nil
}

func validateStruct(typeSpec *ast.TypeSpec) error {
    if structType, ok := typeSpec.Type.(*ast.StructType); ok {
        for _, field := range structType.Fields.List {
            if field.Tag != nil {
                tags := field.Tag.Value
                if err := validateFieldTags(field.Names[0].Name, tags); err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

func validateFieldTags(fieldName, tags string) error {
    // 检查json标签是否符合camelCase
    // 检查db标签是否符合snake_case
    // 验证命名转换是否正确
    
    jsonTag := extractTag(tags, "json")
    dbTag := extractTag(tags, "db")
    
    expectedJSON := toCanelCase(dbTag)
    if jsonTag != expectedJSON {
        return fmt.Errorf("字段 %s 的JSON标签不匹配: 期望 %s, 实际 %s", 
                         fieldName, expectedJSON, jsonTag)
    }
    
    return nil
}

// 转换函数
func toCanelCase(snake string) string {
    if snake == "id" {
        return "id"
    }
    
    parts := strings.Split(snake, "_")
    if len(parts) == 1 {
        return parts[0]
    }
    
    result := parts[0]
    for _, part := range parts[1:] {
        result += strings.Title(part)
    }
    return result
}
```

#### 3.3 强制一致性措施

**CI/CD集成检查**:
```yaml
# .github/workflows/field-mapping-check.yml
name: Field Mapping Consistency Check

on: [push, pull_request]

jobs:
  mapping-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.24.6
          
      - name: Run SQLC Generate
        run: sqlc generate
        
      - name: Validate Field Mapping
        run: go run tools/mapping-validator.go
        
      - name: Check Git Diff
        run: |
          if [ -n "$(git diff --exit-code)" ]; then
            echo "❌ 代码生成后有变更，请运行 sqlc generate 并提交"
            git diff
            exit 1
          fi
          echo "✅ 字段映射一致性检查通过"
```

**代码提交前置检查**:
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "🔍 运行字段映射一致性检查..."

# 生成最新代码
sqlc generate
swag init -g cmd/server/main.go
npx openapi-typescript api.yaml -o types/api.ts

# 验证映射一致性
go run tools/mapping-validator.go
if [ $? -ne 0 ]; then
    echo "❌ 字段映射一致性检查失败"
    exit 1
fi

echo "✅ 字段映射一致性检查通过"
```

## 🔗 URL结构规范

### 1. 基础URL结构
```
https://api.ycgcloud.com/api/v1/{resource}[/{id}][/{action}]
```

#### URL组成部分
- **基础域名**: `https://api.ycgcloud.com`
- **API前缀**: `/api`
- **版本号**: `/v1` (主版本.次版本)
- **资源名**: `/users`, `/files`, `/teams` (复数形式)
- **资源ID**: `/{id}` (唯一标识符)
- **操作动作**: `/{action}` (特殊操作)

#### 环境配置
| 环境 | 域名 | 说明 |
|------|------|------|
| 生产环境 | `https://api.ycgcloud.com` | 正式生产环境API |
| 预发布环境 | `https://api-staging.ycgcloud.com` | 预发布测试环境 |
| 测试环境 | `https://api-test.ycgcloud.com` | 功能测试环境 |
| 开发环境 | `http://localhost:8080` | 本地开发环境 |

### 2. 资源命名规范

#### 标准资源路径
```bash
# ✅ 正确的资源路径
/api/v1/users                    # 用户集合
/api/v1/users/123                # 特定用户
/api/v1/files                    # 文件集合
/api/v1/files/456                # 特定文件
/api/v1/teams                    # 团队集合
/api/v1/teams/789                # 特定团队

# 嵌套资源
/api/v1/users/123/files          # 用户的文件列表
/api/v1/teams/789/members        # 团队成员列表
/api/v1/files/456/versions       # 文件版本历史
/api/v1/files/456/shares         # 文件分享记录

# ❌ 错误示例
/api/v1/getUsers                 # 不在URL中使用动词
/api/v1/user                     # 集合资源应使用复数
/api/v1/Users                    # 不使用大写
/api/v1/user_list                # 不使用下划线
/api/v1/file-management          # 避免过于复杂的路径
```

#### 特殊操作路径
```bash
# 认证相关
POST   /api/v1/auth/login        # 用户登录
POST   /api/v1/auth/logout       # 用户登出
POST   /api/v1/auth/refresh      # 刷新Token
POST   /api/v1/auth/register     # 用户注册
POST   /api/v1/auth/forgot-password    # 忘记密码
POST   /api/v1/auth/reset-password     # 重置密码

# 用户管理
GET    /api/v1/users/profile     # 获取用户资料
PUT    /api/v1/users/profile     # 更新用户资料
GET    /api/v1/users/settings    # 获取用户设置
PUT    /api/v1/users/settings    # 更新用户设置
POST   /api/v1/users/avatar      # 上传头像
DELETE /api/v1/users/avatar      # 删除头像
GET    /api/v1/users/sessions    # 获取用户会话列表
DELETE /api/v1/users/sessions/{id}  # 删除指定会话
DELETE /api/v1/users/sessions    # 删除所有会话

# 文件操作
GET    /api/v1/files             # 获取文件列表
GET    /api/v1/files/{id}        # 获取文件详情
PUT    /api/v1/files/{id}        # 更新文件信息
DELETE /api/v1/files/{id}        # 删除文件（移至回收站）
POST   /api/v1/files/upload      # 文件上传
GET    /api/v1/files/{id}/download     # 文件下载
POST   /api/v1/files/{id}/copy         # 复制文件
POST   /api/v1/files/{id}/move         # 移动文件
POST   /api/v1/files/{id}/share        # 分享文件
DELETE /api/v1/files/{id}/share        # 取消分享

# 搜索和统计
GET    /api/v1/search/files      # 搜索文件
GET    /api/v1/search/users      # 搜索用户
GET    /api/v1/stats/dashboard   # 仪表板统计
GET    /api/v1/stats/usage       # 使用量统计

# 批量操作
POST   /api/v1/files/batch       # 批量文件操作
DELETE /api/v1/files/batch       # 批量删除文件
PUT    /api/v1/users/batch       # 批量更新用户
```

## 📋 API版本控制策略

### 1. 版本控制方案选择

**选定策略：URL路径版本控制**
- **选择理由**：
  - 可读性强，URL中直接体现版本信息
  - 便于缓存，不同版本有独立的URL
  - 简单明了，前端开发友好
  - 支持版本间的渐进式迁移

**URL版本格式**：`/api/v{major}[.{minor}]`

```bash
# 主版本号（重大变更，向下不兼容）
/api/v1/users              # 版本 1.x
/api/v2/users              # 版本 2.x（不兼容v1）

# 可选次版本号（向下兼容的功能增强）
/api/v1.1/users            # 版本 1.1（兼容v1.0）
/api/v1.2/users            # 版本 1.2（兼容v1.0、v1.1）
```

### 2. 版本演进规则

#### 2.1 版本号规则 (Semantic Versioning)

```yaml
版本格式: v{major}.{minor}.{patch}

# 主版本号 (Major)：不兼容的API变更
- 移除接口或字段
- 修改响应数据结构
- 更改HTTP状态码含义
- 修改认证方式

# 次版本号 (Minor)：向下兼容的功能增加
- 新增接口
- 新增可选字段
- 新增可选参数
- 新增响应字段

# 修订号 (Patch)：向下兼容的bug修复
- 修复错误响应
- 性能优化
- 文档修正
- 内部逻辑优化
```

#### 2.2 版本兼容性策略

**向下兼容原则**：
```go
// 版本兼容性检查
type APICompatibility struct {
    // 可以安全添加的变更
    SafeChanges []string `json:"safe_changes"`
    // 需要版本升级的变更
    BreakingChanges []string `json:"breaking_changes"`
}

var CompatibilityRules = APICompatibility{
    SafeChanges: []string{
        "添加新的可选字段",
        "添加新的响应字段",
        "添加新的查询参数（可选）",
        "添加新的HTTP头（可选）",
        "添加新的状态码响应",
        "扩展枚举值（向前兼容）",
    },
    BreakingChanges: []string{
        "删除或重命名字段",
        "修改字段类型",
        "修改必填字段为可选",
        "修改可选字段为必填",
        "移除或修改现有状态码",
        "修改错误响应格式",
        "修改认证方式",
    },
}
```

### 3. 版本生命周期管理

#### 3.1 版本状态流转

```mermaid
graph LR
    A[开发中] --> B[内测版]
    B --> C[公测版]
    C --> D[稳定版]
    D --> E[维护版]
    E --> F[废弃版]
    F --> G[停用版]
```

**版本状态说明**：
```yaml
development:     # 开发中
  description: "内部开发阶段，API可能频繁变更"
  access: "仅内部开发团队"
  stability: "不稳定"

alpha:          # 内测版
  description: "核心功能基本完成，限量用户测试"
  access: "邀请制内测用户"
  stability: "基本稳定，可能有小幅调整"

beta:           # 公测版
  description: "功能完整，公开测试阶段"
  access: "所有注册用户"
  stability: "稳定，仅bug修复"

stable:         # 稳定版
  description: "正式发布版本，生产环境推荐"
  access: "所有用户"
  stability: "高度稳定，仅兼容性修复"

maintenance:    # 维护版
  description: "仅提供关键bug修复和安全更新"
  access: "所有用户"
  stability: "仅维护，不再新增功能"

deprecated:     # 废弃版
  description: "标记为废弃，建议迁移到新版本"
  access: "现有用户可继续使用"
  stability: "停止功能开发"
  migration_deadline: "6个月后停用"

retired:        # 停用版
  description: "已停止服务"
  access: "拒绝访问"
  stability: "不可用"
```

#### 3.2 版本支持策略

**同时支持版本数**: 最多3个主版本
**支持时间表**:
```yaml
v1:
  release_date: "2024-01-01"
  stable_until: "2024-12-31"
  maintenance_until: "2025-06-30"
  retirement_date: "2025-12-31"

v2:
  release_date: "2024-06-01"
  stable_until: "2025-06-30"
  maintenance_until: "2026-01-31"
  retirement_date: "2026-06-30"

v3:
  release_date: "2025-01-01"
  stable_until: "2026-01-31"
  # ... 继续演进
```

### 4. 版本迁移指南

#### 4.1 客户端迁移支持

**渐进式迁移**：
```javascript
// 前端版本适配器
class APIVersionAdapter {
    constructor(targetVersion = 'v1') {
        this.version = targetVersion;
        this.baseURL = `/api/${this.version}`;
    }
    
    // 版本兼容性处理
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        
        // 添加版本信息到请求头
        const headers = {
            'API-Version': this.version,
            'Accept': `application/vnd.ycgcloud.${this.version}+json`,
            ...options.headers
        };
        
        return fetch(url, { ...options, headers });
    }
    
    // 自动版本升级检测
    async checkVersionUpdate() {
        const response = await this.request('/version/info');
        const versionInfo = await response.json();
        
        if (versionInfo.latest !== this.version) {
            return {
                needUpdate: true,
                latest: versionInfo.latest,
                migrationGuide: versionInfo.migrationGuide
            };
        }
        
        return { needUpdate: false };
    }
}
```

**后端版本处理**：
```go
// 版本路由处理
func setupVersionRoutes(router *gin.Engine) {
    // v1版本API
    v1 := router.Group("/api/v1")
    v1.Use(VersionMiddleware("v1"))
    {
        v1.GET("/users", v1GetUsers)
        v1.POST("/users", v1CreateUser)
    }
    
    // v2版本API
    v2 := router.Group("/api/v2")
    v2.Use(VersionMiddleware("v2"))
    {
        v2.GET("/users", v2GetUsers)
        v2.POST("/users", v2CreateUser)
    }
    
    // 版本信息接口
    router.GET("/api/version/info", getVersionInfo)
}

// 版本中间件
func VersionMiddleware(version string) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Header("API-Version", version)
        c.Header("API-Supported-Versions", "v1,v2")
        c.Set("api_version", version)
        c.Next()
    })
}

// 版本信息响应
func getVersionInfo(c *gin.Context) {
    info := VersionInfo{
        Current:     "v2.1.0",
        Latest:      "v2.1.0",
        Supported:   []string{"v1", "v2"},
        Deprecated:  []string{},
        Retired:     []string{},
        ChangeLog:   "/docs/changelog",
        MigrationGuide: "/docs/migration",
    }
    
    c.JSON(200, info)
}
```

#### 4.2 兼容性测试

**自动化版本兼容性测试**：
```yaml
# .github/workflows/api-compatibility.yml
name: API Compatibility Test

on:
  pull_request:
    paths: ['api/**']

jobs:
  compatibility-test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        api_version: [v1, v2]
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Test Environment
        run: |
          docker-compose up -d api-${{ matrix.api_version }}
          
      - name: Run Compatibility Tests
        run: |
          # 运行现有API测试用例，确保向下兼容
          npm run test:api -- --version=${{ matrix.api_version }}
          
      - name: Generate Compatibility Report
        run: |
          # 生成兼容性报告
          go run tools/compatibility-checker.go \
            --old-version=${{ matrix.api_version }} \
            --new-version=current \
            --output=compatibility-report.json
```

## 🔀 HTTP方法规范

### 1. 标准CRUD操作映射

| HTTP方法 | 操作类型 | 资源示例 | 幂等性 | 安全性 | 常见状态码 |
|---------|---------|----------|--------|--------|------------|
| `GET` | 查询 | `GET /api/v1/users` | 是 | 是 | 200, 404 |
| `GET` | 查询单个 | `GET /api/v1/users/123` | 是 | 是 | 200, 404 |
| `POST` | 创建 | `POST /api/v1/users` | 否 | 否 | 201, 400, 409 |
| `PUT` | 完整更新 | `PUT /api/v1/users/123` | 是 | 否 | 200, 204, 404 |
| `PATCH` | 部分更新 | `PATCH /api/v1/users/123` | 否 | 否 | 200, 204, 404 |
| `DELETE` | 删除 | `DELETE /api/v1/users/123` | 是 | 否 | 204, 404 |

#### 方法特性说明
- **幂等性**: 多次执行相同请求的结果是否一致
- **安全性**: 是否会修改服务器状态
- **缓存性**: GET方法可缓存，其他方法通常不可缓存

### 2. 云盘系统特殊操作方法

#### 文件操作专用接口
```bash
# 文件管理操作 (POST)
POST   /api/v1/files/{id}/copy          # 复制文件
POST   /api/v1/files/{id}/move          # 移动文件
POST   /api/v1/files/{id}/share         # 创建分享链接
POST   /api/v1/files/{id}/restore       # 从回收站恢复
POST   /api/v1/files/batch              # 批量文件操作
POST   /api/v1/files/upload             # 文件上传

# 文件状态操作 (PUT - 幂等)
PUT    /api/v1/files/{id}/favorite      # 设置/取消收藏
PUT    /api/v1/files/{id}/lock          # 锁定/解锁文件
PUT    /api/v1/files/{id}/public        # 设置公开状态

# 文件关系移除 (DELETE)
DELETE /api/v1/files/{id}/share         # 取消分享
DELETE /api/v1/files/{id}/favorite      # 取消收藏
DELETE /api/v1/files/{id}/tags/{tag}    # 移除标签
```

#### 用户和团队操作
```bash
# 用户操作 (POST)
POST   /api/v1/users/{id}/activate      # 激活用户
POST   /api/v1/users/{id}/reset-password # 重置密码
POST   /api/v1/users/{id}/send-invitation # 发送邀请

# 团队操作 (POST)
POST   /api/v1/teams/{id}/invite        # 邀请成员
POST   /api/v1/teams/{id}/transfer      # 转移所有权

# 状态更新 (PUT)
PUT    /api/v1/users/{id}/status        # 更新用户状态
PUT    /api/v1/teams/{id}/settings      # 更新团队设置

# 关系移除 (DELETE)
DELETE /api/v1/teams/{id}/members/{uid} # 移除团队成员
DELETE /api/v1/users/{id}/sessions      # 清除用户会话
```

### 3. 云盘系统方法选择指南

#### GET方法 - 数据查询
```bash
# ✅ 推荐用法
GET /api/v1/users                 # 获取用户列表
GET /api/v1/users/123             # 获取用户详情
GET /api/v1/files?type=image      # 按类型过滤文件
GET /api/v1/files/{id}/versions   # 获取文件版本历史
GET /api/v1/search/files?q=报告   # 全文搜索文件
GET /api/v1/teams/123/members     # 获取团队成员
GET /api/v1/files/{id}/download   # 文件下载(重定向)
GET /api/v1/stats/dashboard       # 获取统计数据

# ❌ 错误用法
GET /api/v1/users/create          # 创建操作应使用POST
GET /api/v1/files/delete/123      # 删除操作应使用DELETE
GET /api/v1/files/upload          # 上传操作应使用POST
```

#### POST方法 - 创建和操作
```bash
# ✅ 推荐用法
POST /api/v1/users                # 创建新用户
POST /api/v1/auth/login           # 用户登录认证
POST /api/v1/files/upload         # 文件上传
POST /api/v1/files/{id}/copy      # 复制文件(非幂等)
POST /api/v1/files/batch          # 批量文件操作
POST /api/v1/teams/{id}/invite    # 邀请团队成员
POST /api/v1/search/advanced      # 复杂搜索(带请求体)

# ❌ 错误用法
POST /api/v1/users/123            # 更新用户应使用PUT/PATCH
POST /api/v1/users/123/get        # 查询用户应使用GET
POST /api/v1/files/{id}/favorite  # 状态切换应使用PUT
```

#### PUT方法 - 完整更新和状态设置
```bash
# ✅ 推荐用法
PUT /api/v1/users/123             # 完整更新用户信息
PUT /api/v1/files/{id}/metadata   # 更新文件元数据
PUT /api/v1/files/{id}/favorite   # 设置收藏状态(幂等)
PUT /api/v1/teams/{id}/settings   # 更新团队设置
PUT /api/v1/users/{id}/avatar     # 更新用户头像

# ❌ 错误用法
PUT /api/v1/files/{id}/copy       # 复制操作应使用POST
PUT /api/v1/files/upload          # 上传操作应使用POST
```

#### PATCH方法 - 部分更新
```bash
# ✅ 推荐用法
PATCH /api/v1/users/123           # 部分更新用户信息
PATCH /api/v1/files/{id}          # 更新文件部分属性
PATCH /api/v1/teams/{id}          # 部分更新团队信息

# 请求体示例
{
  "real_name": "新姓名",
  "phone": "+86-13900139000"
}
```

#### DELETE方法 - 删除和移除
```bash
# ✅ 推荐用法
DELETE /api/v1/users/123          # 删除用户
DELETE /api/v1/files/123          # 删除文件(移至回收站)
DELETE /api/v1/files/{id}/share   # 取消文件分享
DELETE /api/v1/teams/{id}/members/{uid} # 移除团队成员
DELETE /api/v1/files/{id}/tags/{tag}    # 移除文件标签

# ❌ 错误用法
DELETE /api/v1/files/{id}/copy    # 复制操作应使用POST
DELETE /api/v1/users/{id}/login   # 登出操作应使用POST
```

## 📊 请求格式规范

### 1. Content-Type规范

| Content-Type | 使用场景 | 云盘系统应用示例 | 说明 |
|--------------|----------|------------------|------|
| `application/json` | 标准API请求 | 用户注册、文件元数据更新、团队管理 | 默认格式，UTF-8编码 |
| `multipart/form-data` | 文件上传 | 文件上传、头像上传、批量文件上传 | 支持二进制数据 |
| `application/x-www-form-urlencoded` | 表单提交 | 用户登录、密码重置 | 简单表单数据 |
| `application/octet-stream` | 二进制流 | 大文件分片上传、文件下载 | 原始二进制数据 |
| `text/plain` | 纯文本 | 文件内容预览、日志上传 | 纯文本内容 |

### 2. 请求头规范

#### 🔒 必需请求头
```http
# 认证信息 (除公开接口外)
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 内容类型
Content-Type: application/json; charset=utf-8

# 用户代理 (用于统计和兼容性)
User-Agent: YCGCloud-Web/1.0.0 (Windows NT 10.0; Win64; x64)
```

#### 🔧 推荐请求头
```http
# 请求追踪ID (用于日志关联)
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000

# 客户端版本 (用于兼容性控制)
X-Client-Version: 1.2.3

# 语言偏好
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8

# 客户端IP (代理环境)
X-Forwarded-For: 192.168.1.100

# 设备信息 (用于安全检测)
X-Device-ID: device_12345
X-Platform: web
```

#### 🎯 特殊场景请求头
```http
# 文件上传专用
X-Upload-Session-ID: upload_session_12345
X-Chunk-Index: 0
X-Total-Chunks: 10
X-File-Hash: sha256:abc123...

# 分享访问专用
X-Share-Token: share_token_12345
X-Share-Password: encrypted_password

# 团队协作专用
X-Team-ID: team_12345
X-Workspace-ID: workspace_67890

# API限流相关
X-Rate-Limit-Key: user_12345
X-Priority: high
```

### 3. 查询参数规范

#### 🔍 分页参数 (Pagination)
```bash
# 🎯 基于页码的分页 (推荐用于用户界面)
GET /api/v1/files?page=1&size=20&total=true
# 响应包含: data, pagination: {page, size, total, pages}

# 🚀 基于游标的分页 (推荐用于大数据量和实时数据)
GET /api/v1/files?cursor=eyJpZCI6MTIzLCJjcmVhdGVkX2F0IjoiMjAyNC0wMS0wMVQwMDowMDowMFoifQ&limit=20
# 响应包含: data, pagination: {cursor, has_more, next_cursor}

# 📊 基于偏移量的分页 (简单场景，不推荐大数据量)
GET /api/v1/files?offset=0&limit=20
# 响应包含: data, pagination: {offset, limit, total}

# 云盘系统分页示例
GET /api/v1/files?page=1&size=50&folder_id=123    # 文件夹内文件分页
GET /api/v1/shares?cursor=abc123&limit=10         # 分享记录分页
GET /api/v1/teams/456/members?page=1&size=20      # 团队成员分页
```

#### 🔤 排序参数 (Sorting)
```bash
# 🔤 单字段排序
GET /api/v1/files?sort=created_at&order=desc      # 按创建时间降序
GET /api/v1/files?sort=name&order=asc             # 按文件名升序
GET /api/v1/files?sort=size&order=desc            # 按文件大小降序

# 🔀 多字段排序
GET /api/v1/files?sort=type,size,name&order=asc,desc,asc
# 先按类型升序，再按大小降序，最后按名称升序

# ⚡ 简化排序语法 (推荐)
GET /api/v1/files?sort=-created_at,+name,-size    # 前缀 - 表示降序，+ 表示升序

# 云盘系统排序示例
GET /api/v1/files?sort=-modified_at               # 最近修改的文件
GET /api/v1/shares?sort=-created_at,-access_count # 最新且最热门的分享
GET /api/v1/users?sort=+username                  # 用户名字母序
```

#### 🎯 过滤参数 (Filtering)
```bash
# 🎯 精确匹配
GET /api/v1/files?type=image&status=active&folder_id=123
GET /api/v1/users?role=admin&status=active
GET /api/v1/teams?type=enterprise&billing_status=paid

# 📏 范围过滤
GET /api/v1/files?size_gt=1048576&size_lt=104857600        # 1MB-100MB
GET /api/v1/files?created_after=2024-01-01T00:00:00Z      # 2024年后创建
GET /api/v1/files?modified_between=2024-01-01,2024-12-31  # 时间范围
GET /api/v1/users?storage_used_gt=5368709120              # 使用超过5GB

# 🔍 模糊搜索
GET /api/v1/files?name_like=报告&content_contains=总结
GET /api/v1/users?email_like=@company.com
GET /api/v1/files?path_contains=/项目/文档/

# 📋 包含关系
GET /api/v1/files?tags_in=工作,重要&type_in=pdf,docx
GET /api/v1/files?folder_id_not_in=1,2,3&status_not_in=deleted,hidden
GET /api/v1/users?role_in=admin,manager&team_id_in=1,2,3

# 🔗 关联过滤
GET /api/v1/files?shared=true&public=false               # 已分享但非公开
GET /api/v1/files?has_thumbnail=true&virus_scanned=true  # 有缩略图且已扫描
GET /api/v1/users?has_teams=true&two_factor_enabled=true # 有团队且启用2FA
```

#### ✅ 字段选择参数 (Field Selection)
```bash
# ✅ 选择特定字段 (减少数据传输)
GET /api/v1/users?fields=id,username,email,avatar_url
GET /api/v1/files?fields=id,name,size,type,created_at
GET /api/v1/teams?fields=id,name,member_count,created_at

# ❌ 排除特定字段 (排除敏感信息)
GET /api/v1/users?exclude=password_hash,two_factor_secret,deleted_at
GET /api/v1/files?exclude=content_hash,virus_scan_result

# 🔗 包含关联数据 (减少请求次数)
GET /api/v1/files?include=folder,tags,shares,owner
GET /api/v1/teams?include=members,workspace,owner
GET /api/v1/users?include=teams,storage_stats

# 🌳 展开嵌套数据 (深度关联)
GET /api/v1/teams?expand=members.user.profile,workspace.files.metadata
GET /api/v1/files?expand=folder.parent,shares.recipients
GET /api/v1/users?expand=teams.workspace.files

# 云盘系统字段选择示例
GET /api/v1/files?fields=id,name,size&include=folder      # 文件基本信息+文件夹
GET /api/v1/users?exclude=password_hash&include=teams     # 用户信息+团队(排除密码)
```

#### 🔍 搜索参数 (Search)
```bash
# 🔍 全文搜索
GET /api/v1/search?q=项目报告&type=files&scope=all
GET /api/v1/search?q=张三&type=users&team_id=123

# 🎯 高级搜索
GET /api/v1/search/files?q=报告&type=pdf&size_gt=1048576&created_after=2024-01-01
GET /api/v1/search/files?content=总结&tags=重要&folder_path=/项目/

# 🏷️ 标签搜索
GET /api/v1/files?tags=工作,重要&tag_mode=all            # 包含所有标签
GET /api/v1/files?tags=工作,个人&tag_mode=any            # 包含任一标签
```

### 4. 请求体规范

#### 📝 标准JSON请求示例

##### 用户注册请求
```json
// POST /api/v1/users - 创建用户
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePassword123!",
  "real_name": "张三",
  "phone": "+86-13800138000",
  "language": "zh-CN",
  "timezone": "Asia/Shanghai",
  "email_notifications": true,
  "terms_accepted": true
}
```

##### 文件元数据更新请求
```json
// PUT /api/v1/files/123 - 更新文件信息
{
  "name": "项目报告_最终版.pdf",
  "description": "2024年度项目总结报告",
  "tags": ["工作", "重要", "2024"],
  "folder_id": 123,
  "is_public": false,
  "permissions": {
    "can_download": true,
    "can_preview": true,
    "can_comment": false
  }
}
```

##### 团队创建请求
```json
// POST /api/v1/teams - 创建团队
{
  "name": "产品开发团队",
  "description": "负责产品功能开发和维护",
  "type": "enterprise",
  "settings": {
    "max_members": 50,
    "storage_quota": 107374182400,
    "allow_external_sharing": true,
    "require_approval_for_join": true
  },
  "initial_members": [
    {"email": "member1@company.com", "role": "admin"},
    {"email": "member2@company.com", "role": "member"}
  ]
}
```

##### 用户部分更新请求
```json
// PATCH /api/v1/users/123 - 部分更新用户
{
  "real_name": "张三丰",
  "phone": "+86-13900139000",
  "avatar_url": "https://cdn.example.com/avatars/new.jpg",
  "notification_settings": {
    "email_enabled": true,
    "push_enabled": false
  }
}
```

#### 🔄 批量操作请求

##### 批量文件操作
```json
// POST /api/v1/files/batch - 批量文件操作
{
  "action": "move",
  "file_ids": [123, 456, 789],
  "target_folder_id": 999,
  "options": {
    "overwrite_existing": false,
    "preserve_permissions": true,
    "preserve_timestamps": true,
    "notify_users": true,
    "create_activity_log": true
  }
}
```

##### 批量用户邀请
```json
// POST /api/v1/teams/123/invitations/batch - 批量邀请用户
{
  "invitations": [
    {
      "email": "user1@company.com",
      "role": "member",
      "message": "欢迎加入我们的团队！"
    },
    {
      "email": "user2@company.com",
      "role": "admin",
      "message": "邀请您担任团队管理员"
    }
  ],
  "options": {
    "send_email": true,
    "expire_days": 7,
    "require_approval": false
  }
}
```

##### 批量删除操作
```json
// DELETE /api/v1/files/batch - 批量删除文件
{
  "file_ids": [123, 456, 789],
  "permanent": false,
  "reason": "清理临时文件",
  "notify_affected_users": true,
  "backup_before_delete": true
}
```

#### 👤 用户管理API请求

##### 获取用户资料
```json
// GET /api/v1/users/profile - 获取当前用户资料
// 无请求体，通过Token识别用户
```

##### 更新用户资料
```json
// PUT /api/v1/users/profile - 更新用户资料
{
  "real_name": "张三",
  "phone": "13800138000",
  "language": "zh-CN",
  "timezone": "Asia/Shanghai",
  "theme": "dark",
  "email_notifications": true,
  "sms_notifications": false
}
```

##### 获取用户设置
```json
// GET /api/v1/users/settings - 获取用户设置
// 无请求体，返回用户的个性化设置
```

##### 更新用户设置
```json
// PUT /api/v1/users/settings - 更新用户设置
{
  "sync_enabled": true,
  "sync_bandwidth_limit": 10485760,
  "auto_preview": true,
  "thumbnail_quality": "high",
  "video_quality": "1080p",
  "auto_cleanup_trash": true,
  "trash_retention_days": 30,
  "share_default_expire": 7,
  "share_require_password": false,
  "notification_file_shared": true,
  "notification_file_commented": true
}
```

##### 上传用户头像
```bash
// POST /api/v1/users/avatar - 上传头像
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary

------WebKitFormBoundary
Content-Disposition: form-data; name="avatar"; filename="avatar.jpg"
Content-Type: image/jpeg

[图片二进制数据]
------WebKitFormBoundary--
```

##### 获取用户会话
```json
// GET /api/v1/users/sessions - 获取用户所有会话
// 无请求体，返回当前用户的所有活跃会话
```

#### 📁 文件上传请求

##### 单文件上传 (multipart/form-data)
```bash
POST /api/v1/files/upload
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="file"; filename="项目报告.pdf"
Content-Type: application/pdf

[PDF文件二进制数据]
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="folder_id"

123
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="tags"

工作,重要
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="description"

2024年度项目总结报告
------WebKitFormBoundary7MA4YWxkTrZu0gW--
```

##### 分片上传初始化
```json
// POST /api/v1/files/upload/init - 分片上传初始化
{
  "filename": "大文件.zip",
  "file_size": 1073741824,
  "file_type": "application/zip",
  "md5_hash": "d41d8cd98f00b204e9800998ecf8427e",
  "sha256_hash": "abc123def456...",
  "parent_id": 456,
  "chunk_size": 2097152,
  "total_chunks": 512,
  "description": "项目资源包",
  "tags": ["项目", "资源"]
}
```

##### 分片上传块
```json
// POST /api/v1/files/upload/chunk - 上传文件块
{
  "upload_session_id": "upload_session_12345",
  "chunk_index": 0,
  "chunk_hash": "sha256:chunk_hash_here",
  "is_last_chunk": false
}
```

##### 批量文件上传
```json
// POST /api/v1/files/upload/batch - 批量文件上传
{
  "folder_id": 123,
  "files": [
    {
      "name": "文档1.pdf",
      "size": 1048576,
      "hash": "sha256:abc123...",
      "type": "application/pdf",
      "description": "重要文档"
    },
    {
      "name": "图片1.jpg",
      "size": 2097152,
      "hash": "sha256:def456...",
      "type": "image/jpeg",
      "description": "产品截图"
    }
  ],
  "options": {
    "overwrite_existing": false,
    "virus_scan": true,
    "generate_thumbnail": true,
    "extract_metadata": true
  }
}
```

## 📤 响应格式规范

### 1. 统一响应结构

```json
// 成功响应格式
{
  "success": true,
  "code": 200,
  "message": "操作成功",
  "data": {
    // 实际数据内容
  },
  "meta": {
    // 元数据信息
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z",
    "version": "1.0.0"
  }
}

// 错误响应格式
{
  "success": false,
  "code": 400,
  "message": "请求参数错误",
  "error": {
    "type": "VALIDATION_ERROR",
    "details": [
      {
        "field": "email",
        "code": "INVALID_FORMAT",
        "message": "邮箱格式不正确"
      }
    ]
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z",
    "version": "1.0.0"
  }
}
```

### 2. 数据响应格式

#### 单个资源响应
```json
// GET /api/v1/users/123
{
  "success": true,
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "real_name": "张三",
    "avatar_url": "https://cdn.example.com/avatars/123.jpg",
    "role": "user",
    "status": "active",
    "storage_quota": 10737418240,
    "storage_used": 5368709120,
    "storage_usage_percentage": 50.0,
    "created_at": "2024-11-01T09:00:00Z",
    "updated_at": "2024-12-01T10:00:00Z"
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z"
  }
}
```

#### 列表资源响应
```json
// GET /api/v1/files?page=1&size=20
{
  "success": true,
  "code": 200,
  "message": "获取成功",
  "data": {
    "items": [
      {
        "id": 456,
        "filename": "项目计划书.docx",
        "file_size": 2048000,
        "file_type": "document",
        "created_at": "2024-11-15T14:30:00Z"
      },
      {
        "id": 457,
        "filename": "设计图.png",
        "file_size": 1024000,
        "file_type": "image",
        "created_at": "2024-11-16T10:15:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 20,
      "total": 156,
      "total_pages": 8,
      "has_previous": false,
      "has_next": true,
      "previous_page": null,
      "next_page": 2
    }
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z",
    "execution_time": 45
  }
}
```

#### 创建资源响应
```json
// POST /api/v1/users
{
  "success": true,
  "code": 201,
  "message": "创建成功",
  "data": {
    "id": 789,
    "username": "new_user",
    "email": "new@example.com",
    "status": "pending",
    "created_at": "2024-12-01T10:00:00Z"
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z",
    "location": "/api/v1/users/789"
  }
}
```

#### 批量操作响应
```json
// POST /api/v1/files/batch
{
  "success": true,
  "code": 200,
  "message": "批量操作完成",
  "data": {
    "total": 10,
    "successful": 8,
    "failed": 2,
    "results": [
      {
        "file_id": 123,
        "status": "success",
        "message": "移动成功"
      },
      {
        "file_id": 124,
        "status": "error",
        "message": "权限不足",
        "error_code": "AUTH_PERMISSION_DENIED"
      }
    ]
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z",
    "execution_time": 1250
  }
}
```

### 3. 特殊响应格式

#### 文件上传响应
```json
// POST /api/v1/files/upload
{
  "success": true,
  "code": 200,
  "message": "上传初始化成功",
  "data": {
    "upload_id": "upload_abcd1234",
    "chunk_size": 2097152,
    "total_chunks": 512,
    "upload_urls": [
      {
        "chunk_number": 1,
        "upload_url": "https://oss.example.com/upload/chunk1?signature=xxx",
        "expires_at": "2024-12-01T11:00:00Z"
      }
    ],
    "upload_token": "token_xyz789",
    "expires_at": "2024-12-01T11:00:00Z"
  }
}
```

#### 搜索结果响应
```json
// GET /api/v1/search/files?q=项目
{
  "success": true,
  "code": 200,
  "message": "搜索成功",
  "data": {
    "query": "项目",
    "total": 25,
    "took": 45,
    "items": [
      {
        "id": 456,
        "filename": "项目计划书.docx",
        "highlights": {
          "filename": ["<em>项目</em>计划书.docx"],
          "description": ["2024年度<em>项目</em>规划"]
        },
        "score": 0.95
      }
    ],
    "suggestions": ["项目管理", "项目计划", "项目总结"],
    "facets": {
      "file_type": {
        "document": 15,
        "image": 8,
        "video": 2
      },
      "created_year": {
        "2024": 20,
        "2023": 5
      }
    }
  }
}
```

## ❌ 错误处理规范

### 1. HTTP状态码规范

| 状态码 | 含义 | 使用场景 |
|-------|------|----------|
| `200` | OK | 成功处理请求 |
| `201` | Created | 成功创建资源 |
| `204` | No Content | 成功处理但无返回内容 |
| `400` | Bad Request | 请求参数错误 |
| `401` | Unauthorized | 未认证或认证失败 |
| `403` | Forbidden | 已认证但权限不足 |
| `404` | Not Found | 资源不存在 |
| `409` | Conflict | 资源冲突 |
| `422` | Unprocessable Entity | 请求格式正确但语义错误 |
| `429` | Too Many Requests | 请求频率限制 |
| `500` | Internal Server Error | 服务器内部错误 |
| `502` | Bad Gateway | 网关错误 |
| `503` | Service Unavailable | 服务不可用 |

### 2. 错误响应格式

#### 标准错误响应
```json
{
  "success": false,
  "code": 400,
  "message": "请求参数错误",
  "error": {
    "type": "VALIDATION_ERROR",
    "code": "INVALID_PARAMETERS",
    "details": [
      {
        "field": "email",
        "code": "INVALID_FORMAT",
        "message": "邮箱格式不正确",
        "value": "invalid-email"
      },
      {
        "field": "password",
        "code": "TOO_SHORT",
        "message": "密码长度至少8位",
        "min_length": 8
      }
    ]
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z",
    "documentation_url": "https://docs.ycgcloud.com/errors#VALIDATION_ERROR"
  }
}
```

#### 认证错误响应
```json
{
  "success": false,
  "code": 401,
  "message": "认证失败",
  "error": {
    "type": "AUTHENTICATION_ERROR",
    "code": "TOKEN_EXPIRED",
    "message": "访问令牌已过期",
    "expires_at": "2024-12-01T09:30:00Z"
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z",
    "refresh_endpoint": "/api/v1/auth/refresh"
  }
}
```

#### 权限错误响应
```json
{
  "success": false,
  "code": 403,
  "message": "权限不足",
  "error": {
    "type": "PERMISSION_ERROR",
    "code": "INSUFFICIENT_PERMISSIONS",
    "message": "您没有删除此文件的权限",
    "required_permission": "file:delete",
    "current_permissions": ["file:read", "file:write"]
  }
}
```

#### Meta字段场景化补充原则

`meta`字段根据不同错误类型提供针对性的补充信息：

| 错误类型 | Meta补充字段 | 用途说明 |
|---------|-------------|----------|
| **验证错误** | `documentation_url` | 指向详细错误文档，帮助开发者理解具体校验规则 |
| **认证错误** | `refresh_endpoint` | 提供令牌刷新地址，便于客户端自动重新认证 |
| **权限错误** | `required_permissions` | 说明所需权限，指导权限申请流程 |
| **限流错误** | `retry_after`, `quota_info` | 重试建议时间和配额使用情况 |
| **资源错误** | `alternative_endpoints` | 推荐的替代API端点 |
| **系统错误** | `incident_id`, `status_page` | 故障跟踪ID和系统状态页面 |

**设计原则**：
- 基础字段（`request_id`, `timestamp`）所有响应必含
- 场景字段按错误类型动态添加，提供实用的错误恢复指导
- 避免敏感信息泄露，仅提供客户端处理所需的最小信息集

#### 资源不存在错误
```json
{
  "success": false,
  "code": 404,
  "message": "资源不存在",
  "error": {
    "type": "RESOURCE_ERROR",
    "code": "FILE_NOT_FOUND",
    "message": "文件不存在或已被删除",
    "resource_type": "file",
    "resource_id": "123"
  }
}
```

#### 业务逻辑错误
```json
{
  "success": false,
  "code": 422,
  "message": "操作失败",
  "error": {
    "type": "BUSINESS",
    "code": "USER_BUSINESS_STORAGE_QUOTA_EXCEEDED",
    "message": "存储空间不足，无法上传文件",
    "details": {
      "current_usage": 10737418240,
      "quota_limit": 10737418240,
      "required_space": 1048576
    }
  }
}
```

#### meta字段的场景化补充原则

meta字段用于提供错误相关的元数据和补充信息，根据不同错误类型包含相应的辅助字段：

**基础字段（所有错误响应必须包含）：**
```json
{
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-12-01T10:00:00Z"
  }
}
```

**场景化补充字段：**

1. **标准错误响应**：包含`documentation_url`
   ```json
   {
     "documentation_url": "https://docs.ycgcloud.com/errors#VALIDATION_ERROR"
   }
   ```

2. **认证错误响应**：包含`refresh_endpoint`
   ```json
   {
     "refresh_endpoint": "/api/v1/auth/refresh"
   }
   ```

3. **权限错误响应**：包含`permission_request_url`
   ```json
   {
     "permission_request_url": "/api/v1/permissions/request"
   }
   ```

4. **资源错误响应**：包含`search_endpoint`
   ```json
   {
     "search_endpoint": "/api/v1/search"
   }
   ```

5. **业务逻辑错误响应**：包含`help_center_url`
   ```json
   {
     "help_center_url": "https://help.ycgcloud.com/storage-quota"
   }
   ```

6. **系统错误响应**：包含`status_page_url`
   ```json
   {
     "status_page_url": "https://status.ycgcloud.com"
   }
   ```

### 3. 云盘系统错误码设计规范

#### 错误码分类体系
```
错误码格式: {业务模块}_{错误类型}_{具体错误}

业务模块:
- AUTH: 认证授权
- USER: 用户管理
- FILE: 文件管理
- FOLDER: 文件夹管理
- SHARE: 分享管理
- TEAM: 团队协作
- UPLOAD: 文件上传
- STORAGE: 存储管理
- SEARCH: 搜索功能
- SYSTEM: 系统级别

错误类型:
- VALIDATION: 验证错误
- PERMISSION: 权限错误
- RESOURCE: 资源错误
- BUSINESS: 业务逻辑错误
- SYSTEM: 系统错误

注：采用三级结构提供更精确的错误分类和处理策略
```

#### 常用错误码定义

##### 认证授权错误 (AUTH_*)
```json
{
  "AUTH_VALIDATION_TOKEN_EXPIRED": "访问令牌已过期",
  "AUTH_VALIDATION_TOKEN_INVALID": "访问令牌无效",
  "AUTH_VALIDATION_REFRESH_TOKEN_EXPIRED": "刷新令牌已过期",
  "AUTH_BUSINESS_LOGIN_FAILED": "用户名或密码错误",
  "AUTH_BUSINESS_ACCOUNT_LOCKED": "账户已被锁定",
  "AUTH_BUSINESS_TWO_FACTOR_REQUIRED": "需要双因子认证",
  "AUTH_VALIDATION_TWO_FACTOR_INVALID": "双因子认证码无效",
  "AUTH_PERMISSION_ACCESS_DENIED": "权限不足",
  "AUTH_BUSINESS_SESSION_EXPIRED": "会话已过期"
}
```

##### 用户管理错误 (USER_*)
```json
{
  "USER_RESOURCE_NOT_FOUND": "用户不存在",
  "USER_VALIDATION_EMAIL_EXISTS": "邮箱已被注册",
  "USER_VALIDATION_USERNAME_EXISTS": "用户名已被占用",
  "USER_VALIDATION_EMAIL_NOT_VERIFIED": "邮箱未验证",
  "USER_BUSINESS_ACCOUNT_DISABLED": "账户已被禁用",
  "USER_BUSINESS_STORAGE_QUOTA_EXCEEDED": "存储配额已满",
  "USER_BUSINESS_BANDWIDTH_LIMIT_EXCEEDED": "带宽限制已达上限",
  "USER_VALIDATION_PASSWORD_TOO_WEAK": "密码强度不足",
  "USER_SYSTEM_PROFILE_UPDATE_FAILED": "用户资料更新失败"
}
```

##### 文件管理错误 (FILE_*)
```json
{
  "FILE_RESOURCE_NOT_FOUND": "文件不存在",
  "FILE_PERMISSION_ACCESS_DENIED": "文件访问被拒绝",
  "FILE_VALIDATION_SIZE_TOO_LARGE": "文件大小超过限制",
  "FILE_VALIDATION_TYPE_NOT_ALLOWED": "文件类型不被允许",
  "FILE_VALIDATION_NAME_INVALID": "文件名包含非法字符",
  "FILE_BUSINESS_ALREADY_EXISTS": "文件已存在",
  "FILE_BUSINESS_VIRUS_DETECTED": "检测到病毒，文件被拒绝",
  "FILE_BUSINESS_CORRUPTED": "文件已损坏",
  "FILE_SYSTEM_PROCESSING_FAILED": "文件处理失败",
  "FILE_SYSTEM_DOWNLOAD_FAILED": "文件下载失败"
}
```

##### 文件上传错误 (UPLOAD_*)
```json
{
  "UPLOAD_SESSION_EXPIRED": "上传会话已过期",
  "UPLOAD_CHUNK_MISSING": "文件块缺失",
  "UPLOAD_CHUNK_INVALID": "文件块校验失败",
  "UPLOAD_SIZE_MISMATCH": "文件大小不匹配",
  "UPLOAD_HASH_MISMATCH": "文件哈希值不匹配",
  "UPLOAD_CONCURRENT_LIMIT": "并发上传数量超限",
  "UPLOAD_BANDWIDTH_EXCEEDED": "上传带宽超限",
  "UPLOAD_STORAGE_FULL": "存储空间不足"
}
```

##### 分享管理错误 (SHARE_*)
```json
{
  "SHARE_NOT_FOUND": "分享链接不存在",
  "SHARE_EXPIRED": "分享链接已过期",
  "SHARE_PASSWORD_REQUIRED": "需要分享密码",
  "SHARE_PASSWORD_INCORRECT": "分享密码错误",
  "SHARE_ACCESS_LIMIT_EXCEEDED": "分享访问次数超限",
  "SHARE_PERMISSION_DENIED": "分享权限不足",
  "SHARE_DISABLED": "分享已被禁用"
}
```

##### 团队协作错误 (TEAM_*)
```json
{
  "TEAM_NOT_FOUND": "团队不存在",
  "TEAM_MEMBER_LIMIT_EXCEEDED": "团队成员数量超限",
  "TEAM_INVITATION_EXPIRED": "团队邀请已过期",
  "TEAM_ALREADY_MEMBER": "用户已是团队成员",
  "TEAM_PERMISSION_DENIED": "团队权限不足",
  "TEAM_STORAGE_QUOTA_EXCEEDED": "团队存储配额已满",
  "TEAM_WORKSPACE_ACCESS_DENIED": "工作空间访问被拒绝"
}
```

#### 错误码分类
```
{业务模块}_{错误类型}_{具体错误}

业务模块:
- AUTH: 认证授权
- USER: 用户管理
- FILE: 文件管理
- FOLDER: 文件夹管理
- SHARE: 分享管理
- TEAM: 团队协作
- UPLOAD: 文件上传
- STORAGE: 存储管理
- SEARCH: 搜索功能
- SYSTEM: 系统级别

错误类型:
- VALIDATION: 验证错误
- PERMISSION: 权限错误
- RESOURCE: 资源错误
- BUSINESS: 业务逻辑错误
- SYSTEM: 系统错误
```

#### 常用错误码定义
```json
{
  "AUTH_VALIDATION_INVALID_CREDENTIALS": "用户名或密码错误",
  "AUTH_PERMISSION_TOKEN_EXPIRED": "访问令牌已过期",
  "AUTH_PERMISSION_TOKEN_INVALID": "访问令牌无效",
  
  "USER_VALIDATION_EMAIL_REQUIRED": "邮箱地址不能为空",
  "USER_VALIDATION_EMAIL_FORMAT": "邮箱格式不正确",
  "USER_RESOURCE_NOT_FOUND": "用户不存在",
  "USER_BUSINESS_ALREADY_EXISTS": "用户已存在",
  
  "FILE_VALIDATION_FILENAME_REQUIRED": "文件名不能为空",
  "FILE_VALIDATION_FILE_TOO_LARGE": "文件大小超出限制",
  "FILE_RESOURCE_NOT_FOUND": "文件不存在",
  "FILE_PERMISSION_ACCESS_DENIED": "文件访问被拒绝",
  "FILE_BUSINESS_STORAGE_QUOTA_EXCEEDED": "存储空间不足",
  
  "TEAM_VALIDATION_NAME_REQUIRED": "团队名称不能为空",
  "TEAM_RESOURCE_NOT_FOUND": "团队不存在",
  "TEAM_PERMISSION_MEMBER_ONLY": "仅团队成员可访问",
  "TEAM_BUSINESS_MEMBER_LIMIT_EXCEEDED": "团队成员数量已达上限"
}
```

## 🔐 认证授权规范

### 1. 认证方式

#### JWT Token认证机制

##### Token结构设计
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT",
    "kid": "key_id_123"
  },
  "payload": {
    "user_id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "permissions": ["file:read", "file:write", "share:create"],
    "team_ids": [1, 2, 3],
    "storage_quota": 10737418240,
    "exp": 1704110400,
    "iat": 1704106800,
    "iss": "ycgcloud.com",
    "aud": "ycgcloud-api",
    "jti": "token_unique_id_123"
  }
}
```

##### Token使用规范
```http
# 标准认证头
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 可选的设备标识
X-Device-ID: device_12345
X-Client-Version: 1.0.0
```

##### Token生命周期管理
```json
{
  "access_token": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "token_type": "Bearer"
  },
  "refresh_token": {
    "token": "rt_abc123def456...",
    "expires_in": 2592000
  },
  "scope": "read write delete"
}
```

##### Token刷新机制
```bash
# Token刷新请求
POST /api/v1/auth/refresh
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "refresh_token": "rt_abc123def456...",
  "device_id": "device_12345"
}

# 刷新响应
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "rt_new_token...",
    "expires_in": 3600,
    "token_type": "Bearer"
  }
}
```

#### API Key认证 (用于服务间调用)
```bash
# 请求头格式
X-API-Key: sk-live_1234567890abcdef
X-API-Secret: secret_abcdef1234567890
X-Service-Name: file-processor
```

#### 双因子认证 (2FA)
```bash
# 启用2FA的登录流程
POST /api/v1/auth/login
{
  "username": "john_doe",
  "password": "password123",
  "totp_code": "123456",
  "device_trust": true
}

# 2FA验证
POST /api/v1/auth/verify-2fa
{
  "temp_token": "temp_token_123",
  "totp_code": "123456",
  "backup_code": "backup_code_789"
}
```

### 2. 权限控制

#### 基于角色的权限控制 (RBAC)

##### 系统角色定义
```json
{
  "roles": {
    "super_admin": {
      "name": "超级管理员",
      "description": "系统最高权限",
      "permissions": ["*"]
    },
    "admin": {
      "name": "管理员",
      "description": "组织管理权限",
      "permissions": [
        "user:*", "team:*", "storage:manage", 
        "system:monitor", "audit:read"
      ]
    },
    "team_owner": {
      "name": "团队所有者",
      "description": "团队完全控制权限",
      "permissions": [
        "team:manage", "member:*", "workspace:*",
        "file:*", "share:*"
      ]
    },
    "team_admin": {
      "name": "团队管理员",
      "description": "团队管理权限",
      "permissions": [
        "member:invite", "member:remove", "workspace:manage",
        "file:*", "share:create"
      ]
    },
    "team_member": {
      "name": "团队成员",
      "description": "团队基础权限",
      "permissions": [
        "file:read", "file:write", "file:upload",
        "share:create", "comment:*"
      ]
    },
    "user": {
      "name": "普通用户",
      "description": "个人空间权限",
      "permissions": [
        "file:personal", "share:personal", "profile:manage"
      ]
    },
    "guest": {
      "name": "访客",
      "description": "只读访问权限",
      "permissions": ["file:read", "share:access"]
    }
  }
}
```

##### 用户权限信息结构
```json
{
  "user_id": 123,
  "username": "john_doe",
  "roles": ["user", "team_admin"],
  "global_permissions": [
    "file:read", "file:write", "file:delete",
    "team:manage_members", "share:create"
  ],
  "resource_permissions": {
    "files": {
      "456": {
        "permissions": ["read", "write", "delete"],
        "inherited_from": "owner",
        "expires_at": null
      },
      "789": {
        "permissions": ["read", "comment"],
        "inherited_from": "team_member",
        "expires_at": "2024-12-31T23:59:59Z"
      }
    },
    "teams": {
      "101": {
        "role": "admin",
        "permissions": ["manage_members", "manage_workspace"],
        "joined_at": "2024-01-01T00:00:00Z"
      }
    },
    "folders": {
      "folder_123": {
        "permissions": ["read", "write", "create_subfolder"],
        "inherited_from": "parent_folder"
      }
    }
  },
  "temporary_permissions": [
    {
      "resource_type": "file",
      "resource_id": "temp_file_456",
      "permissions": ["read", "download"],
      "expires_at": "2024-01-15T12:00:00Z",
      "granted_by": "user_789"
    }
  ]
}
```

#### 权限检查示例

##### 文件操作权限验证
```bash
# 读取文件 - 需要读取权限
GET /api/v1/files/456
# 权限要求: 
# - 全局权限: file:read 或 file:*
# - 资源权限: 对文件456的read权限
# - 继承权限: 父文件夹的read权限
# - 分享权限: 通过分享链接访问

# 编辑文件 - 需要写入权限
PUT /api/v1/files/456
# 权限要求:
# - 全局权限: file:write 或 file:*
# - 资源权限: 对文件456的write权限
# - 文件状态: 文件未被锁定
# - 协作权限: 在协作模式下需要editor角色

# 删除文件 - 需要删除权限
DELETE /api/v1/files/456
# 权限要求:
# - 全局权限: file:delete 或 file:*
# - 资源权限: 对文件456的delete权限
# - 所有者权限: 文件所有者或管理员
# - 团队权限: 团队管理员权限

# 分享文件 - 需要分享权限
POST /api/v1/files/456/shares
# 权限要求:
# - 全局权限: share:create 或 share:*
# - 资源权限: 对文件456的share权限
# - 组织策略: 符合组织分享策略
```

##### 团队操作权限验证
```bash
# 邀请团队成员
POST /api/v1/teams/101/members
# 权限要求:
# - 全局权限: member:invite
# - 团队权限: 对团队101的admin或owner权限
# - 成员限制: 未超过团队成员上限

# 移除团队成员
DELETE /api/v1/teams/101/members/456
# 权限要求:
# - 全局权限: member:remove
# - 团队权限: 对团队101的admin或owner权限
# - 层级限制: 不能移除更高权限的成员

# 管理工作空间
PUT /api/v1/teams/101/workspaces/789
# 权限要求:
# - 全局权限: workspace:manage
# - 团队权限: 对团队101的admin权限
# - 工作空间权限: 对工作空间789的管理权限
```

##### 权限验证响应示例
```json
// 权限不足响应
{
  "success": false,
  "code": "AUTH_PERMISSION_DENIED",
  "message": "权限不足",
  "error": {
    "type": "PermissionDenied",
    "details": {
      "required_permissions": ["file:delete"],
      "user_permissions": ["file:read", "file:write"],
      "resource_id": "file_456",
      "resource_type": "file",
      "action": "delete"
    },
    "suggestions": [
      "请联系文件所有者获取删除权限",
      "或联系团队管理员提升权限"
    ]
  }
}

// 权限检查成功响应
{
  "success": true,
  "data": {
    "permission_granted": true,
    "permission_source": "resource_permission",
    "expires_at": null,
    "additional_permissions": ["file:share", "file:comment"]
  }
}
```

### 3. 安全增强机制

#### 请求签名验证

##### 签名算法
```bash
# HMAC-SHA256签名验证（高安全级别API）
X-Signature: sha256=abc123def456...
X-Timestamp: 1701417600
X-Nonce: random_string_12345
X-API-Key: api_key_123

# 签名计算方法
# signature = HMAC-SHA256(api_secret, method + url + timestamp + nonce + body_hash)
```

##### 签名验证示例
```json
{
  "signature_config": {
    "algorithm": "HMAC-SHA256",
    "header_name": "X-Signature",
    "timestamp_tolerance": 300,
    "nonce_cache_duration": 3600,
    "required_for": [
      "admin_operations",
      "file_upload",
      "sensitive_data_access"
    ]
  },
  "signature_verification": {
    "step1": "验证时间戳在容忍范围内",
    "step2": "检查nonce是否已使用",
    "step3": "重新计算签名并比较",
    "step4": "验证API密钥有效性"
  }
}
```

#### IP白名单与地理位置控制

##### IP访问控制
```json
{
  "ip_whitelist_config": {
    "user_id": 123,
    "allowed_ips": [
      "192.168.1.0/24",
      "10.0.0.100",
      "203.0.113.0/24"
    ],
    "blocked_ips": [
      "198.51.100.0/24"
    ],
    "admin_operations_only": true,
    "geo_restrictions": {
      "allowed_countries": ["CN", "US", "JP"],
      "blocked_countries": ["XX"],
      "require_vpn_detection": true
    }
  }
}
```

##### 异常访问检测
```json
{
  "anomaly_detection": {
    "unusual_location": {
      "enabled": true,
      "action": "require_2fa",
      "notification": true
    },
    "multiple_devices": {
      "max_concurrent_sessions": 5,
      "action": "terminate_oldest"
    },
    "suspicious_activity": {
      "rapid_requests": {
        "threshold": 100,
        "window": 60,
        "action": "temporary_block"
      },
      "failed_auth_attempts": {
        "threshold": 5,
        "window": 300,
        "action": "account_lock"
      }
    }
  }
}
```

#### API限流与防护

##### 限流策略
```json
{
  "rate_limiting": {
    "global_limits": {
      "requests_per_minute": 1000,
      "requests_per_hour": 10000,
      "requests_per_day": 100000
    },
    "user_limits": {
      "free_user": {
        "requests_per_minute": 60,
        "upload_bandwidth": "10MB/s",
        "download_bandwidth": "50MB/s"
      },
      "premium_user": {
        "requests_per_minute": 300,
        "upload_bandwidth": "100MB/s",
        "download_bandwidth": "200MB/s"
      },
      "enterprise_user": {
        "requests_per_minute": 1000,
        "upload_bandwidth": "1GB/s",
        "download_bandwidth": "2GB/s"
      }
    },
    "endpoint_specific_limits": {
      "/api/v1/auth/login": {
        "requests_per_minute": 10,
        "burst_allowance": 5
      },
      "/api/v1/files/upload": {
        "concurrent_uploads": 3,
        "max_file_size": "5GB"
      },
      "/api/v1/search": {
        "requests_per_minute": 30,
        "complex_query_limit": 10
      }
    }
  }
}
```

##### 限流响应头
```http
# 限流信息响应头
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1701417660
X-RateLimit-Retry-After: 60

# 限流超出响应
HTTP/1.1 429 Too Many Requests
Retry-After: 60
{
  "success": false,
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "请求频率超限",
  "error": {
    "type": "RateLimitExceeded",
    "details": {
      "limit": 1000,
      "window": 60,
      "retry_after": 60,
      "limit_type": "user_requests_per_minute"
    }
  }
}
```

## 📏 版本控制规范

### 1. 版本策略

#### URL版本控制 (推荐)
```bash
# 主版本控制
/api/v1/users          # 版本1 - 稳定版本
/api/v2/users          # 版本2 - 新功能版本
/api/v3/users          # 版本3 - 重大更新版本

# 小版本控制
/api/v1.1/users        # v1的小版本更新
/api/v2.3/users        # v2的第3个小版本

# 预览版本
/api/beta/users        # Beta测试版本
/api/alpha/users       # Alpha测试版本
```

#### Header版本控制
```bash
# 标准版本头
API-Version: v1
Accept: application/vnd.ycgcloud.v1+json

# 详细版本信息
API-Version: v2.1
Accept: application/vnd.ycgcloud.v2.1+json

# 向前兼容请求
API-Version: v1
API-Compatibility: v2
```

#### 版本生命周期管理
```json
{
  "version_lifecycle": {
    "v1": {
      "status": "deprecated",
      "release_date": "2023-01-01",
      "deprecation_date": "2024-06-01",
      "sunset_date": "2024-12-31",
      "migration_guide": "/docs/migration/v1-to-v2"
    },
    "v2": {
      "status": "stable",
      "release_date": "2024-01-01",
      "features": ["enhanced_search", "batch_operations", "real_time_sync"]
    },
    "v3": {
      "status": "beta",
      "release_date": "2024-10-01",
      "features": ["ai_integration", "advanced_analytics", "multi_cloud_sync"]
    }
  }
}
```

### 2. 版本兼容性策略

#### 向后兼容原则
```json
// v1 响应格式 - 基础版本
{
  "id": 123,
  "name": "张三",
  "email": "zhang@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}

// v2 响应格式 - 扩展字段，保持v1兼容
{
  "id": 123,
  "name": "张三",
  "email": "zhang@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "profile": {
    "avatar": "https://cdn.example.com/avatar.jpg",
    "bio": "云盘用户",
    "preferences": {
      "language": "zh-CN",
      "timezone": "Asia/Shanghai"
    }
  },
  "storage": {
    "quota": 10737418240,
    "used": 5368709120,
    "available": 5368709120
  }
}

// v3 响应格式 - 结构优化，提供兼容性映射
{
  "user": {
    "id": 123,
    "personal_info": {
      "display_name": "张三",
      "email_address": "zhang@example.com",
      "profile_image": "https://cdn.example.com/avatar.jpg"
    },
    "account_info": {
      "created_timestamp": "2024-01-01T00:00:00Z",
      "last_login": "2024-01-15T10:30:00Z",
      "status": "active"
    },
    "storage_info": {
      "total_quota_bytes": 10737418240,
      "used_bytes": 5368709120,
      "available_bytes": 5368709120,
      "quota_type": "premium"
    }
  },
  "_compatibility": {
    "v1_mapping": {
      "id": "user.id",
      "name": "user.personal_info.display_name",
      "email": "user.personal_info.email_address",
      "created_at": "user.account_info.created_timestamp"
    },
    "v2_mapping": {
      "profile.avatar": "user.personal_info.profile_image",
      "storage.quota": "user.storage_info.total_quota_bytes"
    }
  }
}

// v2 响应格式 (向后兼容)
{
  "id": 123,
  "name": "张三",          // 保留原字段
  "full_name": "张三",     // 新增字段
  "email": "zhang@example.com",
  "contact": {            // 新增嵌套结构
    "email": "zhang@example.com",
    "phone": "+86-13800138000"
  }
}
```

#### 废弃字段处理
```json
{
  "id": 123,
  "name": "张三",
  "email": "zhang@example.com",
  "old_field": "废弃值",    // 标记为废弃但继续返回
  "_deprecated_fields": [
    {
      "field": "old_field",
      "reason": "使用new_field替代",
      "removal_version": "v3",
      "removal_date": "2025-06-01"
    }
  ]
}
```

## 🚀 性能优化规范

### 1. 分页优化

#### 偏移分页 (适用于小数据集)
```bash
GET /api/v1/files?page=1&size=20&sort=-created_at
```

#### 游标分页 (适用于大数据集)
```bash
GET /api/v1/files?cursor=eyJpZCI6MTIzLCJjcmVhdGVkX2F0IjoiMjAyNC0xMi0wMVQxMDowMDowMFoifQ&limit=20
```

#### 深度分页优化
```json
{
  "data": {
    "items": [...],
    "pagination": {
      "current_cursor": "eyJpZCI6MTIz...",
      "next_cursor": "eyJpZCI6MTQ0...",
      "has_more": true,
      "estimated_total": 10000,  // 估算总数，避免精确计数
      "limit": 20
    }
  }
}
```

### 2. 字段选择优化

```bash
# 选择特定字段
GET /api/v1/users?fields=id,username,email

# 排除字段
GET /api/v1/users?exclude=created_at,updated_at

# 嵌套字段选择
GET /api/v1/files?fields=id,filename,user.id,user.username
```

### 3. 缓存策略

#### HTTP缓存头
```bash
# 公共资源缓存
Cache-Control: public, max-age=3600

# 私有资源缓存
Cache-Control: private, max-age=300

# 不缓存敏感数据
Cache-Control: no-cache, no-store, must-revalidate

# 条件请求
If-None-Match: "etag-value"
If-Modified-Since: Wed, 01 Dec 2024 10:00:00 GMT
```

#### 响应缓存头
```json
{
  "data": {...},
  "meta": {
    "cache": {
      "etag": "\"abc123def456\"",
      "last_modified": "2024-12-01T10:00:00Z",
      "max_age": 300,
      "cache_key": "user:123:profile"
    }
  }
}
```

### 4. 批量操作优化

```bash
# 批量获取
GET /api/v1/users?ids=123,456,789

# 批量创建
POST /api/v1/users/batch
{
  "users": [
    {"username": "user1", "email": "user1@example.com"},
    {"username": "user2", "email": "user2@example.com"}
  ]
}

# 批量更新
PUT /api/v1/users/batch
{
  "updates": [
    {"id": 123, "username": "new_user1"},
    {"id": 456, "email": "new_email@example.com"}
  ]
}

# 批量删除
DELETE /api/v1/users/batch
{
  "ids": [123, 456, 789]
}
```

## 📚 API文档规范

### 1. OpenAPI规范

#### 基础信息配置
```yaml
openapi: 3.0.3
info:
  title: 网络云盘系统 API
  description: |
    网络云盘系统的RESTful API接口文档
    
    ## 🔐 认证方式
    - **JWT Bearer Token**: 用户身份认证
    - **API Key**: 服务间调用认证
    - **OAuth 2.0**: 第三方应用授权
    
    ## 🚦 限流策略
    - **免费用户**: 60 requests/min
    - **高级用户**: 300 requests/min
    - **企业用户**: 1000 requests/min
    
    ## 📊 监控与状态
    - **系统状态**: [status.ycgcloud.com](https://status.ycgcloud.com)
    - **API监控**: [monitor.ycgcloud.com](https://monitor.ycgcloud.com)
    
  version: 2.1.0
  contact:
    name: YCG Cloud API 支持团队
    email: api-support@ycgcloud.com
    url: https://docs.ycgcloud.com/support
  license:
    name: MIT License
    url: https://opensource.org/licenses/MIT
  termsOfService: https://ycgcloud.com/terms

servers:
  - url: https://api.ycgcloud.com/v2
    description: 🌐 生产环境
  - url: https://api-staging.ycgcloud.com/v2
    description: 🧪 预发布环境
  - url: https://api-dev.ycgcloud.com/v2
    description: 🔧 开发环境
  - url: https://api-beta.ycgcloud.com/v2
    description: 🚀 Beta测试环境
```

#### 安全配置
```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: |
        JWT Bearer Token认证
        
        获取方式:
        ```bash
        POST /api/v2/auth/login
        {
          "username": "your_username",
          "password": "your_password"
        }
        ```
    
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
      description: |
        API Key认证，用于服务间调用
        
        申请方式: 联系管理员或在控制台生成
    
    OAuth2:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://auth.ycgcloud.com/oauth/authorize
          tokenUrl: https://auth.ycgcloud.com/oauth/token
          scopes:
            read: 读取权限
            write: 写入权限
            admin: 管理权限

security:
  - BearerAuth: []
  - ApiKeyAuth: []
```

### 2. 接口文档模板

#### 标准接口文档
```yaml
paths:
  /files/{file_id}:
    get:
      summary: 📄 获取文件信息
      description: |
        根据文件ID获取文件的详细信息，包括元数据、权限、版本历史等。
        
        ### 🔒 权限要求
        - 需要对文件的 `file:read` 权限
        - 或者通过有效的分享链接访问
        - 团队成员自动拥有团队文件的读取权限
        
        ### 🚦 限流规则
        - 普通用户: 60次/分钟
        - 高级用户: 300次/分钟
        - 企业用户: 1000次/分钟
        
        ### 💡 使用建议
        - 使用 `include` 参数按需获取关联数据
        - 大文件建议先获取元数据再决定是否下载
        - 支持 ETag 缓存，减少不必要的请求
        
      tags:
        - 📁 文件管理
      operationId: getFileById
      parameters:
        - name: file_id
          in: path
          required: true
          description: 文件唯一标识符
          schema:
            type: string
            pattern: '^file_[a-zA-Z0-9_-]+$'
            example: "file_abc123def456"
        - name: include
          in: query
          description: 包含的关联数据
          schema:
            type: array
            items:
              type: string
              enum: [metadata, permissions, versions, comments, shares]
            example: ["metadata", "permissions"]
        - name: version
          in: query
          description: 指定文件版本
          schema:
            type: string
            example: "v1.2.0"
      responses:
        '200':
          description: ✅ 文件信息获取成功
          headers:
            ETag:
              description: 文件版本标识
              schema:
                type: string
                example: '"abc123def456"'
            Last-Modified:
              description: 最后修改时间
              schema:
                type: string
                format: date-time
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/FileResponse'
              examples:
                basic_file:
                  summary: 基础文件信息
                  value:
                    success: true
                    data:
                      id: "file_abc123def456"
                      name: "项目文档.pdf"
                      size: 2048576
                      type: "application/pdf"
                      created_at: "2024-01-01T00:00:00Z"
                      updated_at: "2024-01-15T10:30:00Z"
                with_metadata:
                  summary: 包含元数据的文件信息
                  value:
                    success: true
                    data:
                      id: "file_abc123def456"
                      name: "项目文档.pdf"
                      size: 2048576
                      type: "application/pdf"
                      metadata:
                        author: "张三"
                        title: "项目需求文档"
                        pages: 25
                        created_with: "Microsoft Word"
        '304':
          description: 📋 文件未修改 (基于ETag)
        '403':
          description: 🚫 权限不足
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
              example:
                success: false
                code: "FILE_ACCESS_DENIED"
                message: "文件访问被拒绝"
        '404':
          description: ❌ 文件不存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
              example:
                success: false
                code: "FILE_NOT_FOUND"
                message: "文件不存在"
        '429':
          description: 🚦 请求频率超限
          headers:
            Retry-After:
              description: 重试等待时间(秒)
              schema:
                type: integer
                example: 60
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
      security:
        - BearerAuth: []
        - ApiKeyAuth: []
```

### 3. 多语言SDK示例

#### JavaScript/TypeScript
```typescript
// 安装: npm install @ycgcloud/api-client
import { YCGCloudClient, FileInclude } from '@ycgcloud/api-client';

const client = new YCGCloudClient({
  apiKey: process.env.YCG_API_KEY,
  baseURL: 'https://api.ycgcloud.com/v2',
  timeout: 30000
});

// 获取文件信息
async function getFileInfo(fileId: string) {
  try {
    const response = await client.files.get(fileId, {
      include: [FileInclude.METADATA, FileInclude.PERMISSIONS]
    });
    
    console.log('文件信息:', response.data);
    return response.data;
  } catch (error) {
    if (error.code === 'FILE_NOT_FOUND') {
      console.error('文件不存在');
    } else if (error.code === 'FILE_ACCESS_DENIED') {
      console.error('权限不足');
    } else {
      console.error('获取失败:', error.message);
    }
    throw error;
  }
}

// 上传文件
async function uploadFile(file: File, folderId?: string) {
  const uploadSession = await client.files.createUploadSession({
    filename: file.name,
    size: file.size,
    folder_id: folderId
  });
  
  return client.files.uploadChunks(uploadSession.id, file);
}
```

#### Python
```python
# 安装: pip install ycgcloud-sdk
from ycgcloud import YCGCloudClient, FileInclude
from ycgcloud.exceptions import FileNotFoundError, AccessDeniedError

client = YCGCloudClient(
    api_key=os.getenv('YCG_API_KEY'),
    base_url='https://api.ycgcloud.com/v2',
    timeout=30
)

# 获取文件信息
def get_file_info(file_id: str):
    try:
        response = client.files.get(
            file_id=file_id,
            include=[FileInclude.METADATA, FileInclude.PERMISSIONS]
        )
        print(f"文件信息: {response.data}")
        return response.data
    except FileNotFoundError:
        print("文件不存在")
        raise
    except AccessDeniedError:
        print("权限不足")
        raise
    except Exception as e:
        print(f"获取失败: {e}")
        raise

# 批量上传文件
def batch_upload_files(file_paths: list, folder_id: str = None):
    upload_tasks = []
    
    for file_path in file_paths:
        with open(file_path, 'rb') as f:
            upload_session = client.files.create_upload_session(
                filename=os.path.basename(file_path),
                size=os.path.getsize(file_path),
                folder_id=folder_id
            )
            
            task = client.files.upload_file_async(upload_session.id, f)
            upload_tasks.append(task)
    
    # 等待所有上传完成
    results = await asyncio.gather(*upload_tasks)
    return results
```

#### Go
```go
// go get github.com/ycgcloud/go-sdk
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/ycgcloud/go-sdk"
)

func main() {
    client := ycgcloud.NewClient(os.Getenv("YCG_API_KEY"))
    
    // 获取文件信息
    file, err := getFileInfo(client, "file_abc123def456")
    if err != nil {
        fmt.Printf("获取文件失败: %v\n", err)
        return
    }
    
    fmt.Printf("文件信息: %+v\n", file)
}

func getFileInfo(client *ycgcloud.Client, fileID string) (*ycgcloud.File, error) {
    ctx := context.Background()
    
    file, err := client.Files.Get(ctx, fileID, &ycgcloud.FileGetOptions{
        Include: []string{"metadata", "permissions"},
    })
    
    if err != nil {
        switch {
        case ycgcloud.IsNotFoundError(err):
            return nil, fmt.Errorf("文件不存在: %w", err)
        case ycgcloud.IsAccessDeniedError(err):
            return nil, fmt.Errorf("权限不足: %w", err)
        default:
            return nil, fmt.Errorf("获取失败: %w", err)
        }
    }
    
    return file, nil
}

// 文件上传示例
func uploadFile(client *ycgcloud.Client, filePath string, folderID *string) error {
    ctx := context.Background()
    
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("打开文件失败: %w", err)
    }
    defer file.Close()
    
    stat, err := file.Stat()
    if err != nil {
        return fmt.Errorf("获取文件信息失败: %w", err)
    }
    
    // 创建上传会话
    session, err := client.Files.CreateUploadSession(ctx, &ycgcloud.CreateUploadSessionRequest{
        Filename: stat.Name(),
        Size:     stat.Size(),
        FolderID: folderID,
    })
    if err != nil {
        return fmt.Errorf("创建上传会话失败: %w", err)
    }
    
    // 执行上传
    result, err := client.Files.Upload(ctx, session.ID, file)
    if err != nil {
        return fmt.Errorf("上传失败: %w", err)
    }
    
    fmt.Printf("上传成功: %+v\n", result)
    return nil
}
```

## 🎯 最佳实践

### 1. API设计原则

#### RESTful设计
- ✅ 使用名词表示资源，动词表示操作
- ✅ 保持URL简洁和一致性
- ✅ 正确使用HTTP状态码
- ✅ 支持内容协商

#### 错误处理
- ✅ 提供清晰的错误信息
- ✅ 使用标准化的错误码
- ✅ 包含解决建议
- ✅ 记录详细的错误日志

#### 性能优化
- ✅ 实现合理的缓存策略
- ✅ 支持字段选择和过滤
- ✅ 使用分页避免大数据集
- ✅ 提供批量操作接口

### 2. 安全最佳实践

#### 认证授权
- ✅ 使用强密码策略
- ✅ 实现多因子认证
- ✅ 定期轮换API密钥
- ✅ 最小权限原则

#### 数据保护
- ✅ 敏感数据加密传输
- ✅ 不在URL中传递敏感信息
- ✅ 实现数据脱敏
- ✅ 定期安全审计

#### 防护措施
- ✅ 实现限流和熔断
- ✅ 防止SQL注入和XSS
- ✅ 输入验证和输出编码
- ✅ 监控异常访问模式

### 3. 监控与维护

#### 性能监控
- 📊 响应时间监控
- 📊 错误率统计
- 📊 吞吐量分析
- 📊 资源使用情况

#### 日志记录
- 📝 请求/响应日志
- 📝 错误详情记录
- 📝 性能指标日志
- 📝 安全事件日志

#### 版本管理
- 🔄 平滑版本升级
- 🔄 向后兼容保证
- 🔄 废弃通知机制
- 🔄 迁移指导文档

---

## 📞 支持与反馈

### 技术支持
- 📧 **邮箱**: api-support@ycgcloud.com
- 💬 **在线客服**: [support.ycgcloud.com](https://support.ycgcloud.com)
- 📚 **文档中心**: [docs.ycgcloud.com](https://docs.ycgcloud.com)
- 🐛 **问题反馈**: [github.com/ycgcloud/api-issues](https://github.com/ycgcloud/api-issues)

### 社区资源
- 💻 **开发者社区**: [community.ycgcloud.com](https://community.ycgcloud.com)
- 📖 **API更新日志**: [changelog.ycgcloud.com](https://changelog.ycgcloud.com)
- 🎓 **教程和示例**: [examples.ycgcloud.com](https://examples.ycgcloud.com)
- 📊 **系统状态**: [status.ycgcloud.com](https://status.ycgcloud.com)

---

*本文档最后更新时间: 2024年12月15日*  
*文档版本: v2.1.0*  
*API版本: v2.1.0*
{
  "updates": [
    {"id": 123, "status": "active"},
    {"id": 456, "status": "disabled"}
  ]
}
```

## 📚 API文档规范

### 1. OpenAPI文档结构

```yaml
openapi: 3.0.3
info:
  title: YCG Cloud API
  description: 网络云盘系统API接口文档
  version: 1.0.0
  contact:
    name: API Support
    email: api-support@ycgcloud.com
    url: https://docs.ycgcloud.com
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: https://api.ycgcloud.com/api/v1
    description: 生产环境
  - url: https://staging-api.ycgcloud.com/api/v1
    description: 测试环境

paths:
  /users:
    get:
      summary: 获取用户列表
      description: |
        获取系统中的用户列表，支持分页、排序和过滤。
        
        ### 权限要求
        - 需要 `user:read` 权限
        - 管理员可查看所有用户
        - 普通用户只能查看公开信息
        
      parameters:
        - name: page
          in: query
          description: 页码，从1开始
          required: false
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: size
          in: query
          description: 每页数量，最大100
          required: false
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 成功获取用户列表
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserListResponse'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'

components:
  schemas:
    User:
      type: object
      required:
        - id
        - username
        - email
      properties:
        id:
          type: integer
          format: int64
          description: 用户ID
          example: 123
        username:
          type: string
          minLength: 3
          maxLength: 50
          pattern: '^[a-zA-Z0-9_]+$'
          description: 用户名
          example: john_doe
        email:
          type: string
          format: email
          description: 邮箱地址
          example: john@example.com
          
  responses:
    Unauthorized:
      description: 未授权访问
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
            
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

### 2. 接口文档要素

#### 必需要素
- **接口描述**: 清晰说明接口用途
- **请求参数**: 详细参数说明和验证规则
- **响应格式**: 完整的响应结构和示例
- **错误码**: 可能的错误情况和处理方法
- **权限要求**: 调用接口所需的权限
- **限制说明**: 频率限制、大小限制等

#### 示例文档
```markdown
## 创建用户

### 接口信息
- **URL**: `POST /api/v1/users`
- **描述**: 创建新用户账户
- **权限**: 需要 `user:create` 权限

### 请求参数
| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| username | string | 是 | 用户名，3-50位字母数字下划线 | "john_doe" |
| email | string | 是 | 邮箱地址 | "john@example.com" |
| password | string | 是 | 密码，至少8位 | "SecurePassword123!" |
| real_name | string | 否 | 真实姓名 | "张三" |

### 请求示例
```json
{
  "username": "john_doe",
  "email": "john@example.com", 
  "password": "SecurePassword123!",
  "real_name": "张三"
}
```

### 响应示例
```json
{
  "success": true,
  "code": 201,
  "message": "用户创建成功",
  "data": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "status": "pending",
    "created_at": "2024-12-01T10:00:00Z"
  }
}
```

### 错误响应
| 错误码 | HTTP状态 | 说明 |
|--------|----------|------|
| USER_VALIDATION_EMAIL_FORMAT | 400 | 邮箱格式不正确 |
| USER_BUSINESS_ALREADY_EXISTS | 409 | 用户已存在 |
| AUTH_PERMISSION_INSUFFICIENT | 403 | 权限不足 |
```

## 🚦 API限流和监控策略

### 1. 智能限流配置

**基础限流规则：**
```yaml
rate_limits:
  # 全局限流
  global:
    requests_per_second: 1000
    burst_size: 200
    
  # API类别限流
  upload:
    requests_per_minute: 10
    burst: 3
    max_file_size: 10gb
    concurrent_uploads: 3
    
  download:
    requests_per_minute: 30
    burst: 10
    bandwidth_limit: 50mb_per_sec
    
  search:
    requests_per_minute: 60
    burst: 20
    
  auth:
    login_attempts_per_hour: 20
    register_attempts_per_day: 5
    password_reset_per_hour: 3
    
  # 用户等级差异化限流
  user_tiers:
    free:
      daily_api_calls: 1000
      upload_bandwidth: 5mb_per_sec
      download_bandwidth: 10mb_per_sec
      
    vip:
      daily_api_calls: 10000
      upload_bandwidth: 20mb_per_sec
      download_bandwidth: 50mb_per_sec
      
    enterprise:
      daily_api_calls: 100000
      upload_bandwidth: 100mb_per_sec
      download_bandwidth: 200mb_per_sec
```

**动态限流策略：**
```go
// 自适应限流配置
type AdaptiveRateLimiter struct {
    // 基于系统负载的动态调整
    SystemLoadThreshold   float64 `default:"0.8"`    // 系统负载阈值
    MemoryUsageThreshold  float64 `default:"0.85"`   // 内存使用阈值
    
    // 基于用户行为的智能调整
    UserBehaviorAnalysis  bool    `default:"true"`   // 启用用户行为分析
    SuspiciousActivityDetection bool `default:"true"` // 可疑活动检测
    
    // 限流降级策略
    DegradationRules struct {
        HighLoad    int `default:"50"`  // 高负载时限流降级50%
        CriticalLoad int `default:"20"` // 极高负载时限流降级80%
    }
}
```

### 2. API监控和告警

**监控指标配置：**
```yaml
monitoring:
  # 性能指标
  performance_metrics:
    - response_time_p95        # 95分位响应时间
    - response_time_p99        # 99分位响应时间
    - throughput_per_second    # 每秒吞吐量
    - error_rate_percentage    # 错误率百分比
    - cpu_usage_percentage     # CPU使用率
    - memory_usage_percentage  # 内存使用率
    
  # 业务指标
  business_metrics:
    - active_users_count       # 活跃用户数
    - file_upload_count        # 文件上传数量
    - file_download_count      # 文件下载数量
    - storage_usage_bytes      # 存储使用量
    - bandwidth_usage_bytes    # 带宽使用量
    
  # 安全指标
  security_metrics:
    - failed_login_attempts    # 登录失败次数
    - suspicious_activities    # 可疑活动次数
    - blocked_requests_count   # 被阻止的请求数
    - rate_limited_requests    # 被限流的请求数

**告警规则配置：**
```yaml
# Prometheus告警规则配置，用于监控API性能和安全状况
alerts:
  # 性能告警
  performance:
    - name: "API响应时间过高"
      condition: "response_time_p95 > 2000ms"
      severity: "warning"
      duration: "5m"
      
    - name: "API错误率过高"
      condition: "error_rate > 5%"
      severity: "critical"
      duration: "2m"
      
    - name: "服务器负载过高"
      condition: "cpu_usage > 80%"
      severity: "warning"
      duration: "10m"
      
  # 安全告警
  security:
    - name: "暴力破解攻击"
      condition: "failed_login_attempts > 100 in 10m"
      severity: "critical"
      duration: "1m"
      
    - name: "异常大量请求"
      condition: "requests_per_minute > 1000 from single_ip"
      severity: "warning"
      duration: "5m"
```

### 3. API健康检查和熔断

**健康检查配置：**
```yaml
health_checks:
  # 基础健康检查
  basic:
    endpoint: "/health"
    interval: 30s
    timeout: 10s
    
  # 深度健康检查
  deep:
    endpoint: "/health/deep"
    interval: 60s
    timeout: 30s
    checks:
      - database_connection
      - redis_connection
      - oss_storage_access
      - email_service_status
      
  # 业务健康检查
  business:
    endpoint: "/health/business"
    interval: 120s
    checks:
      - file_upload_functionality
      - file_download_functionality
      - user_authentication
      - search_service

# 熔断器配置
circuit_breaker:
  failure_threshold: 5        # 失败阈值
  success_threshold: 3        # 成功阈值
  timeout: 60s               # 熔断超时时间
  max_requests: 10           # 半开状态最大请求数
  
  # 服务依赖熔断
  services:
    database:
      failure_threshold: 3
      timeout: 30s
    oss_storage:
      failure_threshold: 5
      timeout: 60s
    email_service:
      failure_threshold: 10
      timeout: 120s
```

### 4. API版本管理和兼容性

**版本策略：**
```yaml
api_versioning:
  # 版本管理策略
  strategy: "url_path"        # URL路径版本控制
  current_version: "v1"       # 当前版本
  supported_versions: ["v1"]  # 支持的版本列表
  
  # 版本兼容性
  compatibility:
    v1:
      status: "stable"
      deprecation_date: null
      sunset_date: null
      
  # 版本切换策略
  migration:
    gradual_rollout: true     # 渐进式发布
    canary_percentage: 10     # 金丝雀发布比例
    rollback_threshold: 5     # 回滚错误率阈值
```

**API文档自动化：**
```yaml
documentation:
  # 自动生成配置
  auto_generation:
    swagger_enabled: true
    openapi_version: "3.0.3"
    contact_info:
      name: "API团队"
      email: "api@cloudpan.com"
      
  # 文档测试
  testing:
    mock_server_enabled: true
    example_validation: true
    schema_validation: true
```

## 📋 API设计检查清单

### 设计阶段检查
- [ ] URL遵循RESTful规范，使用复数名词表示资源
- [ ] HTTP方法使用正确（GET查询、POST创建、PUT更新、DELETE删除）
- [ ] 状态码使用标准HTTP状态码，含义明确
- [ ] 请求和响应结构设计合理，字段命名一致
- [ ] 错误响应格式统一，包含错误码和详细信息
- [ ] API版本控制策略明确（URL版本或Header版本）
- [ ] 分页、排序、过滤参数设计完整
- [ ] 认证和授权机制设计安全可靠

### 实现阶段检查
- [ ] 参数验证完整，包含必填项和格式校验
- [ ] 错误处理完善，覆盖各种异常情况
- [ ] 字段映射一致性自动化校验（数据库↔Go结构体↔JSON↔前端）
- [ ] 日志记录详细，便于问题排查
- [ ] 性能优化到位，响应时间符合要求
- [ ] 安全防护措施实施（XSS、CSRF、SQL注入等）
- [ ] API文档自动生成和更新
- [ ] 单元测试和集成测试覆盖完整

### 发布阶段检查
- [ ] API向下兼容性验证通过
- [ ] 性能测试达标，按业务场景细分指标：
  - 小文件查询 (< 1MB): 响应时间 < 100ms，并发量 ≥ 2000 QPS
  - 大文件查询 (1-100MB): 响应时间 < 200ms，并发量 ≥ 1000 QPS
  - 文件上传 (< 10MB): 响应时间 < 5s，并发量 ≥ 500 QPS
  - 文件上传 (10MB-1GB): 响应时间 < 30s，并发量 ≥ 100 QPS
  - 用户认证: 响应时间 < 150ms，并发量 ≥ 1500 QPS
  - 搜索操作: 响应时间 < 300ms，并发量 ≥ 800 QPS
  - 批量操作: 响应时间 < 10s，并发量 ≥ 200 QPS
- [ ] 安全测试通过（渗透测试、漏洞扫描）
- [ ] 文档完整准确，示例代码可运行
- [ ] 监控和告警配置完成
- [ ] 发布计划和回滚方案准备就绪