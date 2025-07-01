# 讀取 .env 變數
ENV_FILE := $(if $(wildcard .env),.env,.env.local)

DB_USER     := $(shell grep '^DB_USER=' $(ENV_FILE) | cut -d '=' -f2)
DB_PASSWORD := $(shell grep '^DB_PASSWORD=' $(ENV_FILE) | cut -d '=' -f2)
DB_HOST     := $(shell grep '^DB_HOST=' $(ENV_FILE) | cut -d '=' -f2)
DB_PORT     := $(shell grep '^DB_PORT=' $(ENV_FILE) | cut -d '=' -f2)
DB_NAME     := $(shell grep '^DB_NAME=' $(ENV_FILE) | cut -d '=' -f2)


# MySQL DSN 組合
DB_URL = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
MIGRATIONS = ./db/migrations

# 建立新的 migration 檔案
# 使用範例：
#   make migrate-create name=create_users_table
migrate-create:
	@echo " 建立新 migration 檔案..."
	@if [ -z "$(name)" ]; then \
		echo "❌ 請加上 name 參數：make migrate-create name=create_example_table"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir $(MIGRATIONS) -seq $(name); \
	echo "✅ 檔案已建立於 $(MIGRATIONS)"

# 執行所有尚未執行的 migration
# 使用範例：
#   make migrate-up
migrate-up:
	@echo "⬆️  執行 migrate up..."
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up
	@echo "✅ 資料表已升級完成"

# 回退上一次的 migration（danger!）
# 使用範例：
#   make migrate-down
migrate-down:
	@echo "⬇️  回退上一次 migration..."
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" down
	@echo "✅ 已回退上一版本"

# 砍掉所有資料表（⚠️ 危險操作！僅限測試用）
# 使用範例：
#   make migrate-drop
migrate-drop:
	@echo "💥 全部砍掉中（drop schema）..."
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" drop -f
	@echo "⚠️  所有資料表已被刪除！"

# 強制設定 migration 當前版本
# 使用範例：
#   make migrate-force VERSION=2
migrate-force:
	@echo "🛠️  強制指定 migration 版本為 $(VERSION)..."
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" force $(VERSION)
	@echo "✅ 強制設定完成"

# 整合初始化：drop → up
# 使用範例：
#   make migrate-reset
migrate-reset:
	@echo "Drop ➜ Up 資料庫重新初始化..."
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" drop -f
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up
	@echo "✅ 資料庫重建完成"


# seed資料
# 使用：make db-seed
db-seed:
	@echo "🌱 執行seed..."
	go run cmd/seed/main.go