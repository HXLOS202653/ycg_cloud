#!/bin/bash

# Go代码质量工具安装脚本
# 用于安装和配置所有必要的代码质量检查工具

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Go环境
check_go_environment() {
    log_info "检查Go环境..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go 1.23或更高版本"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "当前Go版本: $GO_VERSION"
    
    # 检查Go版本是否满足要求（1.23+）
    REQUIRED_VERSION="1.23"
    if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
        log_warning "建议使用Go $REQUIRED_VERSION或更高版本"
    fi
    
    # 检查GOPATH和GOBIN
    if [ -z "$GOPATH" ]; then
        export GOPATH=$(go env GOPATH)
        log_info "设置GOPATH: $GOPATH"
    fi
    
    if [ -z "$GOBIN" ]; then
        export GOBIN=$GOPATH/bin
        log_info "设置GOBIN: $GOBIN"
    fi
    
    # 确保GOBIN在PATH中
    if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
        export PATH="$GOBIN:$PATH"
        log_info "将GOBIN添加到PATH"
    fi
    
    log_success "Go环境检查完成"
}

# 安装工具函数
install_tool() {
    local tool_name=$1
    local tool_package=$2
    local tool_command=$3
    
    log_info "检查工具: $tool_name"
    
    if command -v "$tool_command" &> /dev/null; then
        local version=$("$tool_command" --version 2>/dev/null | head -n1 || echo "未知版本")
        log_success "$tool_name 已安装: $version"
        return 0
    fi
    
    log_info "安装 $tool_name..."
    if go install "$tool_package@latest"; then
        log_success "$tool_name 安装成功"
    else
        log_error "$tool_name 安装失败"
        return 1
    fi
}

# 验证工具安装
verify_tool() {
    local tool_name=$1
    local tool_command=$2
    
    if command -v "$tool_command" &> /dev/null; then
        log_success "✓ $tool_name 可用"
        return 0
    else
        log_error "✗ $tool_name 不可用"
        return 1
    fi
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."
    
    directories=(
        "reports"
        "configs"
        "scripts"
        ".git/hooks"
    )
    
    for dir in "${directories[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            log_info "创建目录: $dir"
        fi
    done
    
    log_success "目录创建完成"
}

# 安装所有工具
install_all_tools() {
    log_info "开始安装代码质量工具..."
    
    # 定义工具列表
    declare -A tools=(
        ["golangci-lint"]="github.com/golangci/golangci-lint/cmd/golangci-lint"
        ["gosec"]="github.com/securecodewarrior/gosec/v2/cmd/gosec"
        ["gocyclo"]="github.com/fzipp/gocyclo/cmd/gocyclo"
        ["goimports"]="golang.org/x/tools/cmd/goimports"
        ["govulncheck"]="golang.org/x/vuln/cmd/govulncheck"
        ["staticcheck"]="honnef.co/go/tools/cmd/staticcheck"
        ["ineffassign"]="github.com/gordonklaus/ineffassign"
        ["misspell"]="github.com/client9/misspell/cmd/misspell"
    )
    
    local failed_tools=()
    
    for tool_name in "${!tools[@]}"; do
        if ! install_tool "$tool_name" "${tools[$tool_name]}" "$tool_name"; then
            failed_tools+=("$tool_name")
        fi
    done
    
    if [ ${#failed_tools[@]} -eq 0 ]; then
        log_success "所有工具安装成功"
    else
        log_error "以下工具安装失败: ${failed_tools[*]}"
        return 1
    fi
}

# 验证所有工具
verify_all_tools() {
    log_info "验证工具安装..."
    
    local tools=("golangci-lint" "gosec" "gocyclo" "goimports" "govulncheck" "staticcheck" "ineffassign" "misspell")
    local failed_verifications=()
    
    for tool in "${tools[@]}"; do
        if ! verify_tool "$tool" "$tool"; then
            failed_verifications+=("$tool")
        fi
    done
    
    if [ ${#failed_verifications[@]} -eq 0 ]; then
        log_success "所有工具验证通过"
    else
        log_error "以下工具验证失败: ${failed_verifications[*]}"
        return 1
    fi
}

# 配置Git hooks
setup_git_hooks() {
    log_info "配置Git hooks..."
    
    if [ ! -d ".git" ]; then
        log_warning "当前目录不是Git仓库，跳过Git hooks配置"
        return 0
    fi
    
    # 检查hooks是否存在
    if [ -f ".git/hooks/pre-commit" ]; then
        log_success "pre-commit hook 已存在"
    else
        log_warning "pre-commit hook 不存在，请运行 'make install-hooks' 安装"
    fi
    
    if [ -f ".git/hooks/pre-push" ]; then
        log_success "pre-push hook 已存在"
    else
        log_warning "pre-push hook 不存在，请运行 'make install-hooks' 安装"
    fi
}

# 检查配置文件
check_config_files() {
    log_info "检查配置文件..."
    
    local config_files=(
        ".golangci.yml"
        "configs/gosec.json"
        "configs/gocyclo.yml"
        "Makefile"
    )
    
    for config_file in "${config_files[@]}"; do
        if [ -f "$config_file" ]; then
            log_success "✓ $config_file 存在"
        else
            log_warning "✗ $config_file 不存在"
        fi
    done
}

# 运行快速测试
run_quick_test() {
    log_info "运行快速测试..."
    
    # 测试golangci-lint
    if command -v golangci-lint &> /dev/null; then
        log_info "测试 golangci-lint..."
        if golangci-lint --version > /dev/null 2>&1; then
            log_success "golangci-lint 工作正常"
        else
            log_error "golangci-lint 测试失败"
        fi
    fi
    
    # 测试gosec
    if command -v gosec &> /dev/null; then
        log_info "测试 gosec..."
        if gosec --version > /dev/null 2>&1; then
            log_success "gosec 工作正常"
        else
            log_error "gosec 测试失败"
        fi
    fi
    
    # 测试gocyclo
    if command -v gocyclo &> /dev/null; then
        log_info "测试 gocyclo..."
        if gocyclo --help > /dev/null 2>&1; then
            log_success "gocyclo 工作正常"
        else
            log_error "gocyclo 测试失败"
        fi
    fi
}

# 显示使用说明
show_usage() {
    echo
    log_info "=== 安装完成 ==="
    echo
    echo "现在你可以使用以下命令:"
    echo
    echo "  make lint          # 运行代码检查"
    echo "  make format        # 格式化代码"
    echo "  make security      # 安全扫描"
    echo "  make complexity    # 复杂度检查"
    echo "  make test          # 运行测试"
    echo "  make quality       # 运行所有质量检查"
    echo "  make install-hooks # 安装Git hooks"
    echo
    echo "配置文件位置:"
    echo "  .golangci.yml      # golangci-lint配置"
    echo "  configs/gosec.json # gosec安全扫描配置"
    echo "  configs/gocyclo.yml# gocyclo复杂度配置"
    echo
    echo "报告输出目录: reports/"
    echo
}

# 主函数
main() {
    echo "=== Go代码质量工具安装脚本 ==="
    echo
    
    # 检查是否在项目根目录
    if [ ! -f "go.mod" ]; then
        log_error "请在Go项目根目录运行此脚本（包含go.mod文件的目录）"
        exit 1
    fi
    
    # 执行安装步骤
    check_go_environment
    create_directories
    install_all_tools
    verify_all_tools
    setup_git_hooks
    check_config_files
    run_quick_test
    show_usage
    
    log_success "工具安装脚本执行完成！"
}

# 处理命令行参数
case "${1:-}" in
    --help|-h)
        echo "用法: $0 [选项]"
        echo
        echo "选项:"
        echo "  --help, -h     显示此帮助信息"
        echo "  --verify       仅验证工具安装"
        echo "  --test         仅运行快速测试"
        echo
        exit 0
        ;;
    --verify)
        check_go_environment
        verify_all_tools
        exit $?
        ;;
    --test)
        run_quick_test
        exit $?
        ;;
    "")
        main
        ;;
    *)
        log_error "未知选项: $1"
        echo "使用 --help 查看帮助信息"
        exit 1
        ;;
esac