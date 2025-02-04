# 使用多階段構建
# 建置階段
FROM golang:1.21-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 安裝基本工具
RUN apk add --no-cache git

# 複製 go.mod
COPY go.mod ./

# 由於目前沒有 go.sum，我們可以初始化它
RUN go mod download
RUN go mod tidy

# 複製源代碼
COPY . .

# 編譯應用
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/app

# 最終階段
FROM alpine:latest

WORKDIR /app

# 安裝基本的運行時依賴
RUN apk --no-cache add ca-certificates tzdata

# 複製編譯好的執行檔
COPY --from=builder /app/main .

# 複製靜態文件和配置文件
COPY --from=builder /app/public ./public
COPY --from=builder /app/configs ./configs

# 設置時區
ENV TZ=Asia/Taipei

# 暴露端口（根據你的應用需求修改）
EXPOSE 8080

# 運行應用
CMD ["./main"]