# 个人博客 (PersonalBlog)

基于 Go + Vue 3 的个人博客站点，支持 Markdown 文章发布、后台管理与浏览量统计。

## 技术栈

- **后端**：Go + Gin + GORM + golang-jwt + bcrypt
- **前端**：Vue 3 + TypeScript + Vite + Element Plus + Pinia + Vue Router + axios
- **Markdown**：marked 前端渲染 + highlight.js 代码高亮 + DOMPurify 防 XSS
- **数据库**：SQLite（开发环境，开箱即用）/ MySQL（生产环境）
- **认证**：JWT
- **部署**：前端 Nginx 静态托管 + 后端二进制 + Nginx 反向代理

## 目录结构

```
PersonalBlog/
├── backend/              # Go 后端
│   ├── cmd/server/       # 程序入口
│   ├── internal/
│   │   ├── config/       # 配置加载
│   │   ├── model/        # GORM 模型
│   │   ├── storage/      # DB 初始化/迁移/种子
│   │   ├── middleware/   # JWT 中间件
│   │   ├── handler/      # HTTP 处理器
│   │   └── router/       # 路由注册
│   ├── config.yaml       # 配置文件
│   └── go.mod
├── frontend/             # Vue 3 前端
│   ├── src/
│   │   ├── api/          # axios + API 封装
│   │   ├── components/   # MarkdownEditor
│   │   ├── router/       # 路由 + 守卫
│   │   ├── stores/       # Pinia 认证状态
│   │   ├── styles/       # 全局样式
│   │   └── views/        # blog/ 前台 + admin/ 后台
│   ├── vite.config.ts
│   └── package.json
├── nginx/nginx.conf      # Nginx 部署示例
└── README.md
```

## 快速开始（本地开发）

### 1. 配置后端

开发环境默认使用 **SQLite**，开箱即用、无需安装 MySQL，数据库文件会自动创建在 `backend/data/blog.db`。编辑 `backend/config.yaml` 调整 `jwt.secret` 与 `admin` 账号密码即可：

```yaml
database:
  driver: "sqlite"                                              # 开发默认 sqlite
  dsn: "data/blog.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
jwt:
  secret: "please-change-this-secret-in-production"
  expire_hours: 24
admin:
  username: "admin"
  password: "admin123"
```

> 生产环境改用 MySQL：将 `driver` 改为 `mysql`，填写真实 DSN（见下方「生产环境 MySQL」）。

### 2. 启动后端

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

首次启动会自动建表并写入初始管理员账号（密码以 bcrypt 加密）。控制台显示 `server listening on :8080` 与 `database initialized (driver=sqlite)` 即成功，访问 `http://localhost:8080/health` 返回 `{"status":"ok"}`。

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

打开 `http://localhost:5173`：

- 博客首页：`/`
- 文章详情：`/article/:id`
- 关于：`/about`
- 后台登录：`/admin/login`（使用 config.yaml 中的账号密码登录）

Vite 已配置 `/api` 代理到后端 8080，前后端联调开箱即用。

## 功能说明

**前台**：文章卡片列表 + 分页、Markdown 渲染（代码高亮）、浏览量自动累加。

**后台**（登录后）：
- 仪表盘：文章总数、已发布/草稿数、总浏览量、最近 5 篇文章
- 文章管理：列表、按状态筛选、编辑、删除
- 写文章：Markdown 编辑器（左右分栏实时预览）、摘要/封面/标签、发布或存为草稿

## 生产部署

### 生产环境 MySQL

生产环境使用 MySQL。先创建数据库：

```sql
CREATE DATABASE personal_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

修改 `backend/config.yaml`：

```yaml
database:
  driver: "mysql"
  dsn: "root:yourpassword@tcp(127.0.0.1:3306)/personal_blog?charset=utf8mb4&parseTime=True&loc=Local"
```

首次启动会自动建表并写入管理员账号，迁移逻辑对 SQLite/MySQL 一致。

### 前端构建

```bash
cd frontend
npm run build        # 产物在 frontend/dist
```

### 后端编译

```bash
cd backend
go build -o bin/server ./cmd/server
./bin/server -config config.yaml
```

建议用 `systemd` / `supervisor` / `nssm` 守护后端进程，监听 `127.0.0.1:8080`。

### Nginx 配置

参考 `nginx/nginx.conf`，把 `frontend/dist` 拷贝到服务器（如 `/var/www/blog/dist`），将 `/api/` 反向代理到后端：

```nginx
location / {
    root /var/www/blog/dist;
    try_files $uri $uri/ /index.html;   # SPA history 路由
}
location /api/ {
    proxy_pass http://127.0.0.1:8080;
}
```

```bash
nginx -t && nginx -s reload
```

## 安全提示

- `config.yaml` 中的 `jwt.secret` 与 `admin.password` 请务必在生产环境修改，建议通过环境变量或配置注入覆盖。
- 首次启动后建议登录后台并尽快修改默认密码（当前需通过修改配置后重启重置，后续可增加「修改密码」功能）。
- CORS 当前为开发便利设为 `*`，生产环境建议收紧为你的域名。
