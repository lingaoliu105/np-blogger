# NP-Blogger

NP-Blogger 是一个现代化的博客同步平台，支持将 GitHub 仓库中的文章自动同步到多个技术社区平台。

## 功能特性

- GitHub 账号授权登录
- 自动同步 GitHub 仓库文章
- 支持多平台同步（掘金、CSDN、知乎）
- 可配置的同步间隔
- 文章智能处理和优化
- 支持 Docker 容器化部署
- Kubernetes 集群部署支持

## 环境要求

### 开发环境

- Node.js 18+
- Go 1.21+
- PostgreSQL 14+
- Milvus 向量数据库

### 生产环境

- Docker
- Kubernetes (可选)

## 本地开发

1. 克隆项目
```bash
git clone https://github.com/yourusername/np-blogger.git
cd np-blogger
```

2. 安装依赖
```bash
# 安装前端依赖
npm install

# 安装后端依赖
go mod download
```

3. 配置环境变量

创建 `.env` 文件并配置以下环境变量：
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=np_blogger

# GitHub OAuth 配置
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URL=http://localhost:8080/auth/callback

# RAG 配置
MILVUS_HOST=localhost
MILVUS_PORT=19530
GEMINI_API_KEY=your_gemini_api_key
```

4. 启动开发服务器
```bash
# 启动前端开发服务器
npm run dev

# 启动后端服务器
go run main.go
```

## Docker 部署

1. 构建 Docker 镜像
```bash
docker build -t np-blogger .
```

2. 运行容器
```bash
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your_db_host \
  -e DB_PORT=5432 \
  -e DB_USER=your_db_user \
  -e DB_PASSWORD=your_db_password \
  -e GITHUB_CLIENT_ID=your_github_client_id \
  -e GITHUB_CLIENT_SECRET=your_github_client_secret \
  -e GEMINI_API_KEY=your_gemini_api_key \
  np-blogger
```

## Kubernetes 部署

1. 创建配置文件

更新 `k8s/deployment.yaml` 中的配置：
- 修改镜像名称和标签
- 配置资源限制
- 设置环境变量

2. 创建 Secret
```bash
kubectl create secret generic np-blogger-secrets \
  --from-literal=db_user=your_db_user \
  --from-literal=db_password=your_db_password \
  --from-literal=github_client_id=your_github_client_id \
  --from-literal=github_client_secret=your_github_client_secret \
  --from-literal=gemini_api_key=your_gemini_api_key
```

3. 部署应用
```bash
kubectl apply -f k8s/deployment.yaml
```

## 配置说明

### 数据库配置

- `DB_HOST`: PostgreSQL 数据库主机地址
- `DB_PORT`: 数据库端口（默认：5432）
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称

### GitHub 配置

- `GITHUB_CLIENT_ID`: GitHub OAuth 应用的客户端 ID
- `GITHUB_CLIENT_SECRET`: GitHub OAuth 应用的客户端密钥
- `GITHUB_REDIRECT_URL`: OAuth 回调地址

### RAG 配置

- `MILVUS_HOST`: Milvus 服务器地址
- `MILVUS_PORT`: Milvus 服务器端口
- `GEMINI_API_KEY`: Google Gemini API 密钥

## 常见问题

1. **Q: 如何修改服务器端口？**
   A: 可以通过环境变量 `PORT` 设置，默认为 8080。

2. **Q: 如何配置 GitHub OAuth 应用？**
   A: 在 GitHub 开发者设置中创建新的 OAuth 应用，并确保回调 URL 与配置匹配。

3. **Q: 如何处理数据库连接问题？**
   A: 检查数据库配置是否正确，确保数据库服务器可访问，并且用户具有适当的权限。

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。