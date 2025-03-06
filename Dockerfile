# 使用多阶段构建
FROM node:18-alpine AS frontend
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM golang:1.21-alpine AS backend
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=frontend /app/.next ./.next
COPY --from=frontend /app/public ./public
COPY --from=backend /app/main .
COPY --from=backend /app/config.yaml .

# 设置环境变量
ENV PORT=8080
ENV GIN_MODE=release

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"]