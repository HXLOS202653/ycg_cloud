#!/bin/bash
# 代码格式化脚本
# 使用gofmt和goimports格式化Go代码

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "${BLUE}🎨 Go代码格式化工具${NC}"
echo "${BLUE}==================${NC}"

# 检查是否在Go项目根目录
if [ ! -f "go.mod" ]; then
    echo "${RED}❌ 错误: 未找到 go.mod 文件，请确保在Go项目根目录运行${NC}"
    exit 1
fi

# 检查必要工具是否安装
check_tool() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "${RED}❌ 错误: $1 未安装${NC}"
        echo "${YELLOW}💡 安装方法:${NC}"
        case $1 in
            "gofmt")
                echo "   gofmt 是Go标准工具，应该随Go一起安装"
                ;;
            "goimports")
                echo "   go install golang.org/x/tools/cmd/goimports@latest"
                ;;
        esac
        exit 1
    fi
}

echo "${YELLOW}📋 检查必要工具...${NC}"
check_tool "gofmt"
check_tool "goimports"
echo "${GREEN}✅ 所有工具已安装${NC}"

# 获取所有Go文件
GO_FILES=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -not -path "./.*" | sort)

if [ -z "$GO_FILES" ]; then
    echo "${YELLOW}⚠️  未找到Go文件${NC}"
    exit 0
fi

echo "${BLUE}📁 找到以下Go文件:${NC}"
echo "$GO_FILES" | sed 's/^/  /'
echo

# 解析命令行参数
CHECK_ONLY=false
VERBOSE=false
FIX_IMPORTS=true

while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--check)
            CHECK_ONLY=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        --no-imports)
            FIX_IMPORTS=false
            shift
            ;;
        -h|--help)
            echo "用法: $0 [选项]"
            echo "选项:"
            echo "  -c, --check      只检查格式，不修改文件"
            echo "  -v, --verbose    显示详细输出"
            echo "  --no-imports     跳过import整理"
            echo "  -h, --help       显示此帮助信息"
            exit 0
            ;;
        *)
            echo "${RED}❌ 未知选项: $1${NC}"
            echo "使用 -h 或 --help 查看帮助"
            exit 1
            ;;
    esac
done

# 格式化统计
FORMATTED_COUNT=0
IMPORT_FIXED_COUNT=0
ERROR_COUNT=0

# 1. gofmt 格式化
echo "${YELLOW}🔧 运行 gofmt 格式化...${NC}"
for file in $GO_FILES; do
    if [ "$VERBOSE" = true ]; then
        echo "  检查: $file"
    fi
    
    if [ "$CHECK_ONLY" = true ]; then
        # 只检查，不修改
        if ! gofmt -l "$file" | grep -q "$file"; then
            if [ "$VERBOSE" = true ]; then
                echo "    ${GREEN}✅ 格式正确${NC}"
            fi
        else
            echo "    ${YELLOW}⚠️  需要格式化: $file${NC}"
            FORMATTED_COUNT=$((FORMATTED_COUNT + 1))
        fi
    else
        # 格式化文件
        TEMP_FILE=$(mktemp)
        if gofmt "$file" > "$TEMP_FILE" 2>/dev/null; then
            if ! cmp -s "$file" "$TEMP_FILE"; then
                cp "$TEMP_FILE" "$file"
                echo "    ${GREEN}✅ 已格式化: $file${NC}"
                FORMATTED_COUNT=$((FORMATTED_COUNT + 1))
            elif [ "$VERBOSE" = true ]; then
                echo "    ${GREEN}✅ 无需格式化: $file${NC}"
            fi
        else
            echo "    ${RED}❌ 格式化失败: $file${NC}"
            ERROR_COUNT=$((ERROR_COUNT + 1))
        fi
        rm -f "$TEMP_FILE"
    fi
done

# 2. goimports 整理导入
if [ "$FIX_IMPORTS" = true ]; then
    echo "${YELLOW}📦 运行 goimports 整理导入...${NC}"
    for file in $GO_FILES; do
        if [ "$VERBOSE" = true ]; then
            echo "  检查导入: $file"
        fi
        
        if [ "$CHECK_ONLY" = true ]; then
            # 只检查，不修改
            TEMP_FILE=$(mktemp)
            if goimports "$file" > "$TEMP_FILE" 2>/dev/null; then
                if ! cmp -s "$file" "$TEMP_FILE"; then
                    echo "    ${YELLOW}⚠️  需要整理导入: $file${NC}"
                    IMPORT_FIXED_COUNT=$((IMPORT_FIXED_COUNT + 1))
                elif [ "$VERBOSE" = true ]; then
                    echo "    ${GREEN}✅ 导入正确${NC}"
                fi
            else
                echo "    ${RED}❌ 导入检查失败: $file${NC}"
                ERROR_COUNT=$((ERROR_COUNT + 1))
            fi
            rm -f "$TEMP_FILE"
        else
            # 整理导入
            TEMP_FILE=$(mktemp)
            if goimports "$file" > "$TEMP_FILE" 2>/dev/null; then
                if ! cmp -s "$file" "$TEMP_FILE"; then
                    cp "$TEMP_FILE" "$file"
                    echo "    ${GREEN}✅ 已整理导入: $file${NC}"
                    IMPORT_FIXED_COUNT=$((IMPORT_FIXED_COUNT + 1))
                elif [ "$VERBOSE" = true ]; then
                    echo "    ${GREEN}✅ 无需整理导入: $file${NC}"
                fi
            else
                echo "    ${RED}❌ 导入整理失败: $file${NC}"
                ERROR_COUNT=$((ERROR_COUNT + 1))
            fi
            rm -f "$TEMP_FILE"
        fi
    done
fi

# 显示统计结果
echo
echo "${BLUE}📊 格式化统计:${NC}"
echo "  处理文件数: $(echo "$GO_FILES" | wc -l)"
if [ "$CHECK_ONLY" = true ]; then
    echo "  需要格式化: $FORMATTED_COUNT"
    if [ "$FIX_IMPORTS" = true ]; then
        echo "  需要整理导入: $IMPORT_FIXED_COUNT"
    fi
else
    echo "  已格式化: $FORMATTED_COUNT"
    if [ "$FIX_IMPORTS" = true ]; then
        echo "  已整理导入: $IMPORT_FIXED_COUNT"
    fi
fi
echo "  错误数: $ERROR_COUNT"

# 退出状态
if [ $ERROR_COUNT -gt 0 ]; then
    echo "${RED}❌ 格式化过程中出现错误${NC}"
    exit 1
elif [ "$CHECK_ONLY" = true ] && [ $((FORMATTED_COUNT + IMPORT_FIXED_COUNT)) -gt 0 ]; then
    echo "${YELLOW}⚠️  发现需要格式化的文件${NC}"
    echo "${YELLOW}💡 运行 'make format' 或 './scripts/format.sh' 进行格式化${NC}"
    exit 1
else
    echo "${GREEN}🎉 代码格式化完成！${NC}"
    exit 0
fi