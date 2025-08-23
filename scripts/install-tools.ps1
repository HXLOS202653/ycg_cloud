# Go代码质量工具安装脚本 (PowerShell版本)
# 用于在Windows环境下安装和配置所有必要的代码质量检查工具

param(
    [switch]$Help,
    [switch]$Verify,
    [switch]$Test
)

# 设置错误处理
$ErrorActionPreference = "Stop"

# 颜色定义
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

function Log-Info {
    param([string]$Message)
    Write-ColorOutput "[INFO] $Message" "Blue"
}

function Log-Success {
    param([string]$Message)
    Write-ColorOutput "[SUCCESS] $Message" "Green"
}

function Log-Warning {
    param([string]$Message)
    Write-ColorOutput "[WARNING] $Message" "Yellow"
}

function Log-Error {
    param([string]$Message)
    Write-ColorOutput "[ERROR] $Message" "Red"
}

# 检查Go环境
function Test-GoEnvironment {
    Log-Info "检查Go环境..."
    
    try {
        $goVersion = go version
        if (-not $goVersion) {
            Log-Error "Go未安装，请先安装Go 1.23或更高版本"
            exit 1
        }
        
        $versionMatch = $goVersion -match "go([0-9]+\.[0-9]+)"
        if ($versionMatch) {
            $currentVersion = $Matches[1]
            Log-Info "当前Go版本: $currentVersion"
            
            # 检查版本是否满足要求
            $requiredVersion = [version]"1.23"
            $currentVersionObj = [version]$currentVersion
            
            if ($currentVersionObj -lt $requiredVersion) {
                Log-Warning "建议使用Go 1.23或更高版本"
            }
        }
        
        # 检查和设置Go环境变量
        $goPath = go env GOPATH
        $goBin = go env GOBIN
        
        if (-not $goBin) {
            $goBin = Join-Path $goPath "bin"
        }
        
        Log-Info "GOPATH: $goPath"
        Log-Info "GOBIN: $goBin"
        
        # 检查GOBIN是否在PATH中
        $currentPath = $env:PATH
        if ($currentPath -notlike "*$goBin*") {
            Log-Warning "GOBIN不在PATH中，建议添加 $goBin 到系统PATH"
        }
        
        Log-Success "Go环境检查完成"
        return $true
    }
    catch {
        Log-Error "Go环境检查失败: $($_.Exception.Message)"
        return $false
    }
}

# 安装工具函数
function Install-Tool {
    param(
        [string]$ToolName,
        [string]$ToolPackage,
        [string]$ToolCommand
    )
    
    Log-Info "检查工具: $ToolName"
    
    try {
        $null = Get-Command $ToolCommand -ErrorAction SilentlyContinue
        if ($?) {
            try {
                $version = & $ToolCommand --version 2>$null | Select-Object -First 1
                if (-not $version) {
                    $version = "未知版本"
                }
                Log-Success "$ToolName 已安装: $version"
                return $true
            }
            catch {
                Log-Success "$ToolName 已安装"
                return $true
            }
        }
    }
    catch {
        # 工具未安装，继续安装流程
    }
    
    Log-Info "安装 $ToolName..."
    try {
        $installCmd = "go install $ToolPackage@latest"
        Invoke-Expression $installCmd
        
        if ($LASTEXITCODE -eq 0) {
            Log-Success "$ToolName 安装成功"
            return $true
        }
        else {
            Log-Error "$ToolName 安装失败 (退出码: $LASTEXITCODE)"
            return $false
        }
    }
    catch {
        Log-Error "$ToolName 安装失败: $($_.Exception.Message)"
        return $false
    }
}

# 验证工具安装
function Test-Tool {
    param(
        [string]$ToolName,
        [string]$ToolCommand
    )
    
    try {
        $null = Get-Command $ToolCommand -ErrorAction SilentlyContinue
        if ($?) {
            Log-Success "✓ $ToolName 可用"
            return $true
        }
        else {
            Log-Error "✗ $ToolName 不可用"
            return $false
        }
    }
    catch {
        Log-Error "✗ $ToolName 不可用"
        return $false
    }
}

# 创建必要的目录
function New-RequiredDirectories {
    Log-Info "创建必要的目录..."
    
    $directories = @(
        "reports",
        "configs",
        "scripts",
        ".git\hooks"
    )
    
    foreach ($dir in $directories) {
        if (-not (Test-Path $dir)) {
            try {
                New-Item -ItemType Directory -Path $dir -Force | Out-Null
                Log-Info "创建目录: $dir"
            }
            catch {
                Log-Warning "无法创建目录 $dir : $($_.Exception.Message)"
            }
        }
    }
    
    Log-Success "目录创建完成"
}

# 安装所有工具
function Install-AllTools {
    Log-Info "开始安装代码质量工具..."
    
    # 定义工具列表
    $tools = @{
        "golangci-lint" = "github.com/golangci/golangci-lint/cmd/golangci-lint"
        "gosec" = "github.com/securecodewarrior/gosec/v2/cmd/gosec"
        "gocyclo" = "github.com/fzipp/gocyclo/cmd/gocyclo"
        "goimports" = "golang.org/x/tools/cmd/goimports"
        "govulncheck" = "golang.org/x/vuln/cmd/govulncheck"
        "staticcheck" = "honnef.co/go/tools/cmd/staticcheck"
        "ineffassign" = "github.com/gordonklaus/ineffassign"
        "misspell" = "github.com/client9/misspell/cmd/misspell"
    }
    
    $failedTools = @()
    
    foreach ($toolName in $tools.Keys) {
        $toolPackage = $tools[$toolName]
        if (-not (Install-Tool -ToolName $toolName -ToolPackage $toolPackage -ToolCommand $toolName)) {
            $failedTools += $toolName
        }
    }
    
    if ($failedTools.Count -eq 0) {
        Log-Success "所有工具安装成功"
        return $true
    }
    else {
        Log-Error "以下工具安装失败: $($failedTools -join ', ')"
        return $false
    }
}

# 验证所有工具
function Test-AllTools {
    Log-Info "验证工具安装..."
    
    $tools = @("golangci-lint", "gosec", "gocyclo", "goimports", "govulncheck", "staticcheck", "ineffassign", "misspell")
    $failedVerifications = @()
    
    foreach ($tool in $tools) {
        if (-not (Test-Tool -ToolName $tool -ToolCommand $tool)) {
            $failedVerifications += $tool
        }
    }
    
    if ($failedVerifications.Count -eq 0) {
        Log-Success "所有工具验证通过"
        return $true
    }
    else {
        Log-Error "以下工具验证失败: $($failedVerifications -join ', ')"
        return $false
    }
}

# 配置Git hooks
function Set-GitHooks {
    Log-Info "配置Git hooks..."
    
    if (-not (Test-Path ".git")) {
        Log-Warning "当前目录不是Git仓库，跳过Git hooks配置"
        return
    }
    
    # 检查hooks是否存在
    if (Test-Path ".git\hooks\pre-commit") {
        Log-Success "pre-commit hook 已存在"
    }
    else {
        Log-Warning "pre-commit hook 不存在，请运行 'make install-hooks' 安装"
    }
    
    if (Test-Path ".git\hooks\pre-push") {
        Log-Success "pre-push hook 已存在"
    }
    else {
        Log-Warning "pre-push hook 不存在，请运行 'make install-hooks' 安装"
    }
}

# 检查配置文件
function Test-ConfigFiles {
    Log-Info "检查配置文件..."
    
    $configFiles = @(
        ".golangci.yml",
        "configs\gosec.json",
        "configs\gocyclo.yml",
        "Makefile"
    )
    
    foreach ($configFile in $configFiles) {
        if (Test-Path $configFile) {
            Log-Success "✓ $configFile 存在"
        }
        else {
            Log-Warning "✗ $configFile 不存在"
        }
    }
}

# 运行快速测试
function Invoke-QuickTest {
    Log-Info "运行快速测试..."
    
    # 测试golangci-lint
    try {
        $null = Get-Command "golangci-lint" -ErrorAction SilentlyContinue
        if ($?) {
            Log-Info "测试 golangci-lint..."
            $null = golangci-lint --version 2>$null
            if ($LASTEXITCODE -eq 0) {
                Log-Success "golangci-lint 工作正常"
            }
            else {
                Log-Error "golangci-lint 测试失败"
            }
        }
    }
    catch {
        Log-Error "golangci-lint 测试失败: $($_.Exception.Message)"
    }
    
    # 测试gosec
    try {
        $null = Get-Command "gosec" -ErrorAction SilentlyContinue
        if ($?) {
            Log-Info "测试 gosec..."
            $null = gosec --version 2>$null
            if ($LASTEXITCODE -eq 0) {
                Log-Success "gosec 工作正常"
            }
            else {
                Log-Error "gosec 测试失败"
            }
        }
    }
    catch {
        Log-Error "gosec 测试失败: $($_.Exception.Message)"
    }
    
    # 测试gocyclo
    try {
        $null = Get-Command "gocyclo" -ErrorAction SilentlyContinue
        if ($?) {
            Log-Info "测试 gocyclo..."
            $null = gocyclo --help 2>$null
            if ($LASTEXITCODE -eq 0) {
                Log-Success "gocyclo 工作正常"
            }
            else {
                Log-Error "gocyclo 测试失败"
            }
        }
    }
    catch {
        Log-Error "gocyclo 测试失败: $($_.Exception.Message)"
    }
}

# 显示使用说明
function Show-Usage {
    Write-Host ""
    Log-Info "=== 安装完成 ==="
    Write-Host ""
    Write-Host "现在你可以使用以下命令:"
    Write-Host ""
    Write-Host "  make lint          # 运行代码检查"
    Write-Host "  make format        # 格式化代码"
    Write-Host "  make security      # 安全扫描"
    Write-Host "  make complexity    # 复杂度检查"
    Write-Host "  make test          # 运行测试"
    Write-Host "  make quality       # 运行所有质量检查"
    Write-Host "  make install-hooks # 安装Git hooks"
    Write-Host ""
    Write-Host "配置文件位置:"
    Write-Host "  .golangci.yml      # golangci-lint配置"
    Write-Host "  configs/gosec.json # gosec安全扫描配置"
    Write-Host "  configs/gocyclo.yml# gocyclo复杂度配置"
    Write-Host ""
    Write-Host "报告输出目录: reports/"
    Write-Host ""
}

# 显示帮助信息
function Show-Help {
    Write-Host "Go代码质量工具安装脚本 (PowerShell版本)"
    Write-Host ""
    Write-Host "用法: .\install-tools.ps1 [选项]"
    Write-Host ""
    Write-Host "选项:"
    Write-Host "  -Help              显示此帮助信息"
    Write-Host "  -Verify            仅验证工具安装"
    Write-Host "  -Test              仅运行快速测试"
    Write-Host ""
}

# 主函数
function Main {
    Write-Host "=== Go代码质量工具安装脚本 (PowerShell版本) ===" -ForegroundColor Cyan
    Write-Host ""
    
    # 检查是否在项目根目录
    if (-not (Test-Path "go.mod")) {
        Log-Error "请在Go项目根目录运行此脚本（包含go.mod文件的目录）"
        exit 1
    }
    
    # 执行安装步骤
    if (-not (Test-GoEnvironment)) { exit 1 }
    New-RequiredDirectories
    if (-not (Install-AllTools)) { exit 1 }
    if (-not (Test-AllTools)) { exit 1 }
    Set-GitHooks
    Test-ConfigFiles
    Invoke-QuickTest
    Show-Usage
    
    Log-Success "工具安装脚本执行完成！"
}

# 处理命令行参数
if ($Help) {
    Show-Help
    exit 0
}
elseif ($Verify) {
    Test-GoEnvironment
    Test-AllTools
    exit 0
}
elseif ($Test) {
    Invoke-QuickTest
    exit 0
}
else {
    Main
}