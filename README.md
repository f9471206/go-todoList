# 📝 Go TodoList API

一個使用 Golang + Gin + GORM 開發的 RESTful TodoList API 專案，具備基本的 CRUD、JWT 認證、Swagger 文件、自動化測試與資料遷移支援。

---

## 📦 技術棧

- **Golang** 1.23+
- **Gin** - Web 框架
- **GORM** - ORM 工具，支援 MySQL/PostgreSQL
- **Swagger** - API 文件（/swagger/index.html）
- **JWT** - 登入認證機制
- **Docker** - 可選，支援容器部署
- **sql-migrate** - 資料表遷移管理
- **zap** - 高效能 logger
- **sqlmock + testify** - 單元測試用

---

## 📂 專案結構
```
.
├── cmd/               # 主程式與種子資料入口（如 seed/main.go）
├── config/            # 環境變數與資料庫設定
├── controllers/       # API 控制器（邏輯入口）
├── db/
│   ├── migrations/    # 資料庫 migration 檔案（使用 migrate）
│   └── seed/          # 初始化 seed 資料
├── docs/              # Swagger 文件自動產生（swag init）
├── dto/               # 請求/回應資料轉換物件
├── logs/              # 日誌資料夾（zap logger 輸出）
├── middleware/        # 中介層，JWT/Recovery 等攔截器
├── mocks/             # 使用 gomock 產生的 mock 類別
├── models/            # GORM 資料模型
├── pkg/               # 第三方或輔助工具模組
├── repositories/      # 資料存取層封裝（repository pattern）
├── response/          # 統一 API 回應格式處理
├── routes/            # 路由註冊設定
├── services/          # 業務邏輯層（含單元測試）
├── tmp/               # 開發時暫存目錄（例如 build 暫存）
├── utils/             # 工具方法（例如 logger、加密器）
├── .env.local         # ✅ 本機執行用的環境變數
├── .env.test          # ✅ 測試用環境變數（unit test 專用）
├── .env.docker        # ✅ Docker 環境變數（docker-compose 使用）
├── .env.example       # 📌 範本檔案，供新開發者複製用
├── .fresh.conf        # fresh 工具自動重啟設定
├── Dockerfile         # Docker 映像建構設定
├── docker-compose.yml # 快速啟動整體專案服務
├── go.mod             # Go module 描述
├── go.sum             # 套件版本鎖定
├── main.go            # 入口點：Web Server 啟動點
├── Makefile           # 一鍵建置/遷移/初始化指令整合
└── README.md
```



---

## 🚀 啟動專案

```bash
# 安裝依賴
go mod tidy

# 載入 .env.local 並執行
go run main.go
```

或使用 Docker：
```
docker-compose up --build
```

## 🔄 資料庫操作（migrate）
請先安裝 migrate CLI：

```
brew install golang-migrate
# or
curl -L https://github.com/golang-migrate/migrate/releases/... > /usr/local/bin/migrate
chmod +x /usr/local/bin/migrate
```

## 🌱 種子（seed）
```
make db-seed
```

## ✅ 單元測試
```
go test ./services/...
```

## 🔐 Swagger 文件
```
http://localhost:8080/swagger/index.html
```


## 🧪 生成 Mock

使用 [mockgen](https://github.com/golang/mock) 工具生成介面的 mock 物件：

```
mockgen -source=path/to/your/interface_file.go -destination=path/to/mocks/mock_interface.go -package=mocks
```