# 從 .env 讀變數
DB_USER := $(shell grep '^DB_USER=' .env | cut -d '=' -f2)
DB_PASSWORD := $(shell grep '^DB_PASSWORD=' .env | cut -d '=' -f2)
DB_HOST := $(shell grep '^DB_HOST=' .env | cut -d '=' -f2)
DB_PORT := $(shell grep '^DB_PORT=' .env | cut -d '=' -f2)
DB_NAME := $(shell grep '^DB_NAME=' .env | cut -d '=' -f2)

DB_URL = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
MIGRATIONS=./db/migrations

migrate-up:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" down

migrate-drop:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" drop -f

migrate-force:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" force $(VERSION)
