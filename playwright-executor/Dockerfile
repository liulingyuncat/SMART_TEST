FROM node:20-slim

WORKDIR /app

# 安装依赖
COPY package*.json ./
RUN npm ci --only=production

# 复制源码
COPY server.js ./

# 设置环境变量默认值
ENV PORT=53730
ENV NODE_ENV=production

EXPOSE 53730

CMD ["node", "server.js"]
