#!/bin/bash

# 数据库迁移单元测试执行脚本
# Database Migration Unit Test Runner Script

set -e

echo "🚀 开始执行数据库迁移单元测试..."
echo "Starting database migration unit tests..."

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 设置测试环境变量
export GO_ENV=test
export MIGRATION_TEST_MODE=true

# 检查Go版本
echo "📋 检查Go版本..."
go version

# 检查依赖
echo "📦 检查依赖..."
go mod tidy
go mod verify

# 运行linter检查
echo "🔍 运行代码质量检查..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run ./internal/pkg/migration/...
else
    echo "⚠️  golangci-lint 未安装，跳过linter检查"
fi

# 运行go vet
echo "🔍 运行go vet检查..."
go vet ./internal/pkg/migration/...

# 运行go fmt检查
echo "🔍 检查代码格式..."
if [ -n "$(gofmt -l ./internal/pkg/migration/)" ]; then
    echo "❌ 代码格式不正确，请运行 gofmt -s -w ."
    exit 1
fi

# 创建测试报告目录
mkdir -p test-reports

echo "🧪 运行单元测试..."

# 运行基础测试
echo "📝 运行基础功能测试..."
go test -v ./internal/pkg/migration/ -run "TestMySQLMigrationRecord|TestMongoDBMigrationRecord|TestMigrationConfig|TestValidationConfig|TestBackupConfig|TestMigrationPlan" \
    -coverprofile=test-reports/basic-coverage.out \
    -timeout=30s

# 运行测试工具测试
echo "🛠️  运行测试工具测试..."
go test -v ./internal/pkg/migration/ -run "TestCreateTestMigration|TestGetBasicTestMigrations" \
    -coverprofile=test-reports/utils-coverage.out \
    -timeout=30s

# 运行性能基准测试
echo "⚡ 运行性能基准测试..."
go test -bench=. -benchmem ./internal/pkg/migration/ -run="^$" \
    -timeout=60s > test-reports/benchmark-results.txt

# 生成覆盖率报告
echo "📊 生成测试覆盖率报告..."
go tool cover -func=test-reports/basic-coverage.out > test-reports/coverage-summary.txt
go tool cover -html=test-reports/basic-coverage.out -o test-reports/coverage.html

# 显示覆盖率摘要
echo "📈 测试覆盖率摘要："
cat test-reports/coverage-summary.txt | tail -1

# 检查覆盖率阈值
COVERAGE=$(go tool cover -func=test-reports/basic-coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
THRESHOLD=70

if (( $(echo "$COVERAGE >= $THRESHOLD" | bc -l) )); then
    echo "✅ 测试覆盖率 $COVERAGE% >= $THRESHOLD%，符合要求"
else
    echo "❌ 测试覆盖率 $COVERAGE% < $THRESHOLD%，需要增加测试"
    exit 1
fi

# 运行竞态条件检测
echo "🏃 运行竞态条件检测..."
go test -race ./internal/pkg/migration/ -run="TestCreateTestMigration" -timeout=30s

echo "🎉 所有测试完成！"
echo ""
echo "测试报告位置："
echo "  - 覆盖率HTML报告: test-reports/coverage.html"
echo "  - 覆盖率摘要: test-reports/coverage-summary.txt"
echo "  - 基准测试结果: test-reports/benchmark-results.txt"
echo ""
echo "要查看详细的测试覆盖率报告，请在浏览器中打开: file://$(pwd)/test-reports/coverage.html"