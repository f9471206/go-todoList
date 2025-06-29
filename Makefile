# 讀取 .env 檔案
include .env
export

DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
MIGRATIONS=./db/migrations

migrate-up:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" down

migrate-drop:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" drop -f

migrate-force:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" force $(VERSION)
