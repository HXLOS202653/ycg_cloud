# 安全配置说明

## 环境变量配置

为了保护敏感信息，本项目使用环境变量来管理数据库密码等敏感配置。

### 设置步骤

1. **复制环境变量示例文件**
   ```bash
   cp env.example .env
   ```

2. **编辑 .env 文件，填入真实的敏感信息**
   ```bash
   # 使用你的实际密码替换示例值
   MYSQL_PASSWORD=你的实际MySQL密码
   MONGODB_PASSWORD=你的实际MongoDB密码  
   REDIS_PASSWORD=你的实际Redis密码
   MONGODB_URI=mongodb://root:你的实际MongoDB密码@dbconn.sealosbja.site:43851/?directConnection=true
   ```

3. **确保 .env 文件不被提交到版本控制**
   - `.env` 文件已经添加到 `.gitignore` 中
   - 绝不要将包含真实密码的文件提交到Git仓库

### 配置文件说明

- `configs/development.yaml` - 开发环境配置，使用环境变量
- `configs/production.yaml` - 生产环境配置，使用环境变量  
- `configs/database.yaml` - 数据库专用配置，使用环境变量
- `env.example` - 环境变量示例文件，可以安全提交
- `.env` - 实际环境变量文件，**绝不提交**

### 环境变量格式

配置文件中使用 `${变量名:默认值}` 格式：
- `${MYSQL_PASSWORD}` - 必须设置的变量
- `${MYSQL_HOST:localhost}` - 可选变量，默认值为localhost

### 安全注意事项

⚠️ **重要安全提醒**：
- 绝不要在配置文件中硬编码密码
- 绝不要将 `.env` 文件提交到版本控制
- 定期更换数据库密码
- 在生产环境中使用强密码
- 限制数据库访问权限

### 部署环境配置

不同环境下的环境变量设置：

- **开发环境**: 使用 `.env` 文件
- **CI/CD环境**: 在GitHub Actions中设置Secrets
- **生产环境**: 在服务器或容器编排系统中设置环境变量