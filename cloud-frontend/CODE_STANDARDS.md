# 代码规范文档

## 📋 概述

本项目使用 ESLint + Prettier 来确保代码质量和一致的代码风格。

## 🛠️ 工具配置

### ESLint 9.33.0
- **配置文件**: `eslint.config.js`
- **主要规则**:
  - TypeScript 严格模式
  - React + React Hooks 规则
  - 代码质量检查
  - 无障碍性检查
  - 导入/导出规范

### Prettier 3.6.2
- **配置文件**: `.prettierrc`
- **主要设置**:
  - 单引号优先
  - 无分号
  - 缩进: 2 空格
  - 行宽: 100 字符
  - 尾随逗号: ES5 兼容

### Husky Git Hooks
- **pre-commit**: 自动格式化和代码检查
- **commit-msg**: 提交消息格式验证

## 📜 代码风格规范

### 命名规范
```typescript
// ✅ 变量和函数 - camelCase
const userName = 'john'
const getUserInfo = () => {}

// ✅ 组件和类 - PascalCase
const UserProfile = () => {}
class ApiClient {}

// ✅ 常量 - SCREAMING_SNAKE_CASE
const MAX_FILE_SIZE = 1024

// ✅ 文件名
// 组件文件: PascalCase
UserProfile.tsx
// 工具文件: camelCase
apiClient.ts
// 配置文件: camelCase
designTokens.ts
```

### TypeScript 规范
```typescript
// ✅ 接口定义
interface UserProps {
  id: string
  name: string
  isActive: boolean
}

// ✅ 类型定义
type Status = 'pending' | 'completed' | 'failed'

// ✅ 函数类型
const handleClick: (event: MouseEvent) => void = () => {}
```

### React 规范
```typescript
// ✅ 组件定义
const UserProfile: React.FC<UserProps> = ({ id, name, isActive }) => {
  return <div>{name}</div>
}

// ✅ Hook 使用
const [count, setCount] = useState(0)
const { data, isLoading } = useQuery(['users', id], fetchUser)

// ✅ 事件处理
const handleSubmit = (event: FormEvent) => {
  event.preventDefault()
  // 处理逻辑
}
```

## 🚀 可用命令

### 开发命令
```bash
# 启动开发服务器
npm run dev

# 构建项目
npm run build

# 预览构建结果
npm run preview
```

### 代码质量命令
```bash
# 代码检查
npm run lint

# 自动修复代码问题
npm run lint:fix

# 代码格式化
npm run format

# 检查格式化状态
npm run format:check

# 类型检查
npm run type-check

# 完整代码质量检查
npm run code-quality

# 自动修复 + 格式化
npm run code-fix
```

## 🔧 IDE 配置

### VSCode 设置
项目包含 `.vscode/settings.json` 配置文件，提供：

- 保存时自动格式化
- ESLint 自动修复
- 统一的编辑器设置
- 推荐扩展列表

### 推荐扩展
- **必需**:
  - ESLint
  - Prettier - Code formatter
  - TypeScript and JavaScript Language Features

- **推荐**:
  - GitLens
  - Auto Rename Tag
  - Path Intellisense
  - Thunder Client

## 📝 提交规范

### 提交消息格式
```
type(scope): description

# 类型 (type):
feat:     新功能
fix:      修复bug
docs:     文档更新
style:    代码格式调整
refactor: 代码重构
test:     测试相关
chore:    杂务（配置、依赖等）
perf:     性能优化
ci:       CI/CD相关
build:    构建相关
revert:   回滚

# 示例:
feat(auth): add user login functionality
fix(ui): resolve button alignment issue
docs: update README with setup instructions
```

### Git Hooks 行为
- **pre-commit**: 
  - 运行类型检查
  - 执行 lint-staged（格式化修改的文件）
  - 检查代码格式
  
- **commit-msg**: 
  - 验证提交消息格式

## 🚨 常见问题

### 1. ESLint 错误
```bash
# 查看所有 lint 错误
npm run lint

# 自动修复可修复的错误
npm run lint:fix
```

### 2. 格式化问题
```bash
# 格式化所有文件
npm run format

# 检查哪些文件需要格式化
npm run format:check
```

### 3. TypeScript 错误
```bash
# 类型检查
npm run type-check

# 在开发时，IDE 会实时显示 TypeScript 错误
```

### 4. Git Hook 失败
如果 pre-commit 钩子失败：
1. 修复报错的代码问题
2. 重新添加文件: `git add .`
3. 重新提交: `git commit -m "your message"`

## 🎯 最佳实践

### 1. 开发流程
1. 编写代码
2. 保存文件（VSCode 自动格式化）
3. 运行 `npm run code-quality` 检查
4. 修复问题
5. 提交代码

### 2. 代码审查
- 确保代码通过所有 lint 检查
- 验证类型安全
- 检查命名规范
- 确认代码格式一致

### 3. 团队协作
- 统一使用项目配置
- 不要修改 `.prettierrc` 和 `eslint.config.js` 除非团队同意
- 提交前务必运行 `npm run code-quality`

## 📚 扩展阅读

- [ESLint 官方文档](https://eslint.org/)
- [Prettier 官方文档](https://prettier.io/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [React 官方文档](https://react.dev/)
- [Husky 官方文档](https://typicode.github.io/husky/)

---

*最后更新时间: 2024年12月*