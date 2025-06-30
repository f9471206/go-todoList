FROM golang:1.23

# 安裝 migrate CLI 工具
RUN apt-get update && apt-get install -y curl tar && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate

# 安裝 fresh（即時編譯工具）
RUN go install github.com/pilu/fresh@latest

# 關閉 VCS stamping
ENV GOFLAGS="-buildvcs=false"

WORKDIR /app

# 建立 tmp 目錄並給予可執行權限，避免 permission denied 問題
RUN mkdir -p /app/tmp && chmod -R 777 /app/tmp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["fresh"]
