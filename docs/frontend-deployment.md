# Web 智能测试平台 - 前端部署文档

## 环境要求

- Node.js 14.x 或更高版本
- npm 6.x 或 yarn 1.x

## 安装依赖

```bash
cd frontend
npm install
```

## 开发环境运行

```bash
npm start
```

应用将在 http://localhost:3000 启动

## 构建生产版本

```bash
npm run build
```

构建产物将生成在 `build/` 目录

## 环境配置

### 开发环境 (.env.development)
```
REACT_APP_API_BASE_URL=http://localhost:8080/api/v1
```

### 生产环境 (.env.production)
```
REACT_APP_API_BASE_URL=https://api.production.com/api/v1
```

## 测试

运行所有测试:
```bash
npm test
```

生成覆盖率报告:
```bash
npm test -- --coverage
```

## Nginx 部署配置

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    root /var/www/webtest/frontend/build;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## Docker 部署

### Dockerfile
```dockerfile
FROM node:14-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 构建和运行
```bash
docker build -t webtest-frontend .
docker run -p 80:80 webtest-frontend
```

## 性能优化

1. 代码分割已通过 React.lazy 实现
2. 构建时自动压缩 JS/CSS
3. 使用 gzip 压缩静态资源
4. 配置缓存策略:
   - HTML: no-cache
   - JS/CSS: max-age=31536000

## 监控和日志

建议集成:
- Sentry 错误监控
- Google Analytics 用户分析
- Performance API 性能监控
