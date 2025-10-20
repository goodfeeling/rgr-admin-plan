# Makefile
.PHONY: generate-dict build clean test-connection

# 默认加载 .env 文件
-include .env

# 默认值（从环境变量获取，如果没有则使用默认值）
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432
DB_USER ?= appuser
DB_PASSWORD ?= youShouldChangeThisPassword
DB_NAME ?= boilerplate_go
DB_SSLMODE ?= disable

# 构建 DSN
DICT_DB_DSN ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

generate-dict:
	@echo "Database DSN: $$(echo "$(DICT_DB_DSN)" | sed 's/:[^@]*@/:***@/')"
	DICT_DB_DSN="$(DICT_DB_DSN)" go run cmd/tools/dict_generator/main.go

# 简单的连接测试（手动方式）
test-connection:
	@echo "DSN for manual testing:"
	@echo "$(DICT_DB_DSN)"

generate-dict-build:
	go build -o bin/dict_generator cmd/tools/dict_generator/main.go

generate-dict-run:
	DICT_DB_DSN="$(DICT_DB_DSN)" ./bin/dict_generator

generate-dict-all: generate-dict-build generate-dict-run

build: generate-dict
	go build -o bin/app .

clean:
	rm -f bin/dict_generator