# PowerShell 代码格式化脚本
# 使用gofmt和goimports格式化Go代码

param(
    [switch]$Check,
    [switch]$Verbose,
    [switch]$NoImports,
    [switch]$Help
)

# 颜色函数
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    
    $colorMap = @{
        "Red" = "Red"
        "Green" = "Green"
        "Yellow" = "Yellow"
        "Blue" = "Blue"
        "White" = "White"
    }
    
    Write-Host $Message -ForegroundColor $colorMap[$Color]
}

# 显示帮助信息
if ($Help) {
    Write-ColorOutput "用法: .\scripts\format.ps1 [选项]" "Blue"
    Write-ColorOutput "选项:" "Blue"
    Write-ColorOutput "  -Check       只检查格式，不修改文件" "White"
    Write-ColorOutput "  -Verbose     显示详细输出" "White"
    Write-ColorOutput "  -NoImports   跳过import整理" "White"
    Write-ColorOutput "  -Help        显示此帮助信息" "White"
    exit 0
}

Write-ColorOutput "🎨 Go代码格式化工具" "Blue"
Write-ColorOutput "==================" "Blue"

# 检查是否在Go项目根目录
if (-not (Test-Path "go.mod")) {
    Write-ColorOutput "❌ 错误: 未找到 go.mod 文件，请确保在Go项目根目录运行" "Red"
    exit 1
}

# 检查必要工具是否安装
function Test-Command {
    param([string]$Command)
    
    try {
        Get-Command $Command -ErrorAction Stop | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

Write-ColorOutput "📋 检查必要工具..." "Yellow"

if (-not (Test-Command "gofmt")) {
    Write-ColorOutput "❌ 错误: gofmt 未安装" "Red"
    Write-ColorOutput "💡 gofmt 是Go标准工具，应该随Go一起安装" "Yellow"
    exit 1
}

if (-not (Test-Command "goimports")) {
    Write-ColorOutput "❌ 错误: goimports 未安装" "Red"
    Write-ColorOutput "💡 安装方法: go install golang.org/x/tools/cmd/goimports@latest" "Yellow"
    exit 1
}

Write-ColorOutput "✅ 所有工具已安装" "Green"

# 获取所有Go文件
$goFiles = Get-ChildItem -Path . -Recurse -Filter "*.go" | 
    Where-Object { 
        $_.FullName -notmatch "\\vendor\\" -and 
        $_.FullName -notmatch "\\.git\\" -and
        $_.Name -notmatch "^\."
    } | 
    Sort-Object FullName

if ($goFiles.Count -eq 0) {
    Write-ColorOutput "⚠️  未找到Go文件" "Yellow"
    exit 0
}

Write-ColorOutput "📁 找到以下Go文件:" "Blue"
$goFiles | ForEach-Object { Write-ColorOutput "  $($_.FullName.Replace($PWD.Path, '.'))" "White" }
Write-Host

# 格式化统计
$formattedCount = 0
$importFixedCount = 0
$errorCount = 0

# 1. gofmt 格式化
Write-ColorOutput "🔧 运行 gofmt 格式化..." "Yellow"

foreach ($file in $goFiles) {
    $relativePath = $file.FullName.Replace($PWD.Path, '.')
    
    if ($Verbose) {
        Write-ColorOutput "  检查: $relativePath" "White"
    }
    
    try {
        if ($Check) {
            # 只检查，不修改
            $formatted = & gofmt $file.FullName 2>$null
            $original = Get-Content $file.FullName -Raw
            
            if ($formatted -ne $original) {
                Write-ColorOutput "    ⚠️  需要格式化: $relativePath" "Yellow"
                $formattedCount++
            } elseif ($Verbose) {
                Write-ColorOutput "    ✅ 格式正确" "Green"
            }
        } else {
            # 格式化文件
            $formatted = & gofmt $file.FullName 2>$null
            $original = Get-Content $file.FullName -Raw
            
            if ($formatted -ne $original) {
                $formatted | Set-Content $file.FullName -NoNewline
                Write-ColorOutput "    ✅ 已格式化: $relativePath" "Green"
                $formattedCount++
            } elseif ($Verbose) {
                Write-ColorOutput "    ✅ 无需格式化: $relativePath" "Green"
            }
        }
    }
    catch {
        Write-ColorOutput "    ❌ 格式化失败: $relativePath" "Red"
        $errorCount++
    }
}

# 2. goimports 整理导入
if (-not $NoImports) {
    Write-ColorOutput "📦 运行 goimports 整理导入..." "Yellow"
    
    foreach ($file in $goFiles) {
        $relativePath = $file.FullName.Replace($PWD.Path, '.')
        
        if ($Verbose) {
            Write-ColorOutput "  检查导入: $relativePath" "White"
        }
        
        try {
            if ($Check) {
                # 只检查，不修改
                $formatted = & goimports $file.FullName 2>$null
                $original = Get-Content $file.FullName -Raw
                
                if ($formatted -ne $original) {
                    Write-ColorOutput "    ⚠️  需要整理导入: $relativePath" "Yellow"
                    $importFixedCount++
                } elseif ($Verbose) {
                    Write-ColorOutput "    ✅ 导入正确" "Green"
                }
            } else {
                # 整理导入
                $formatted = & goimports $file.FullName 2>$null
                $original = Get-Content $file.FullName -Raw
                
                if ($formatted -ne $original) {
                    $formatted | Set-Content $file.FullName -NoNewline
                    Write-ColorOutput "    ✅ 已整理导入: $relativePath" "Green"
                    $importFixedCount++
                } elseif ($Verbose) {
                    Write-ColorOutput "    ✅ 无需整理导入: $relativePath" "Green"
                }
            }
        }
        catch {
            Write-ColorOutput "    ❌ 导入整理失败: $relativePath" "Red"
            $errorCount++
        }
    }
}

# 显示统计结果
Write-Host
Write-ColorOutput "📊 格式化统计:" "Blue"
Write-ColorOutput "  处理文件数: $($goFiles.Count)" "White"

if ($Check) {
    Write-ColorOutput "  需要格式化: $formattedCount" "White"
    if (-not $NoImports) {
        Write-ColorOutput "  需要整理导入: $importFixedCount" "White"
    }
} else {
    Write-ColorOutput "  已格式化: $formattedCount" "White"
    if (-not $NoImports) {
        Write-ColorOutput "  已整理导入: $importFixedCount" "White"
    }
}

Write-ColorOutput "  错误数: $errorCount" "White"

# 退出状态
if ($errorCount -gt 0) {
    Write-ColorOutput "❌ 格式化过程中出现错误" "Red"
    exit 1
} elseif ($Check -and ($formattedCount + $importFixedCount) -gt 0) {
    Write-ColorOutput "⚠️  发现需要格式化的文件" "Yellow"
    Write-ColorOutput "💡 运行 'make format' 或 '.\scripts\format.ps1' 进行格式化" "Yellow"
    exit 1
} else {
    Write-ColorOutput "🎉 代码格式化完成！" "Green"
    exit 0
}