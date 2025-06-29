FROM golang:1.23

RUN apt-get update && apt-get install -y curl tar && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate



# 安裝 fresh
RUN go install github.com/pilu/fresh@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["fresh"]
