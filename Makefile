# 数据库连接配置
DB_URL = postgresql://admin:123456@localhost:15432/xledger

# 表与服务目录的映射关系
TABLES = user ledger

# 构建单个表的model (用法: make model-table TABLE=user)
model-table:
	@if [ -z "$(TABLE)" ] ; then \
		exit 1; \
	fi
	@mkdir -p ./service/$(TABLE)/model
	goctl model pg datasource -url="$(DB_URL)" -table="$(TABLE)" -dir="./service/$(TABLE)/model" -cache=true --style=goZero

# 构建ledger服务相关的model
model-all:
	@for table in $(TABLES); do \
		echo "正在生成表: $$table"; \
		mkdir -p ./service/$$table/model && \
		goctl model pg datasource -url="$(DB_URL)" -table="$$table" -dir="./service/$$table/model" -cache=true || true --style=goZero; \
	done

# 构建单个PRC代码
proto-rpc:
	@if [ -z "$(TABLE)" ] ; then \
		exit 1; \
	fi
	goctl rpc protoc service/$(TABLE)/rpc/$(TABLE).proto --style=goZero --go_out=service/$(TABLE)/rpc --go-grpc_out=service/$(TABLE)/rpc --zrpc_out=service/$(TABLE)/rpc

# 构建单个API代码
proto-api:
	@if [ -z "$(TABLE)" ] ; then \
		exit 1; \
	fi
	goctl api go -api service/$(TABLE)/api/$(TABLE).api --style=goZero -dir service/$(TABLE)/api

up:
	@docker-compose up -d

down:
	@docker-compose down

# 数据库迁移文件映射
# 格式: 00001_user.sql, 00002_ledger.sql
MIGRATE_FILES = $(wildcard model/????_*.sql)

# 数据库迁移 (用法: make migrate TABLE=ledger)
migrate:
	@if [ -z "$(TABLE)" ]; then \
		echo "❌ Error: TABLE parameter required. Usage: make migrate TABLE=ledger"; \
		exit 1; \
	fi
	@MIGRATE_FILE=$$(ls model/*_$(TABLE).sql 2>/dev/null | head -n 1); \
	if [ -z "$$MIGRATE_FILE" ]; then \
		echo "❌ Error: Migration file not found for table $(TABLE)"; \
		exit 1; \
	fi; \
	echo "Running migration: $$MIGRATE_FILE"; \
	docker exec -i xledger_postgres psql -U admin -d xledger < $$MIGRATE_FILE && \
	echo "✅ Migration completed for $(TABLE)"

# 数据库回滚 (用法: make migrate-down TABLE=ledger)
migrate-down:
	@if [ -z "$(TABLE)" ]; then \
		echo "❌ Error: TABLE parameter required. Usage: make migrate-down TABLE=ledger"; \
		exit 1; \
	fi
	@echo "⚠️  Rolling back $(TABLE) migration..."
	@if [ "$(TABLE)" = "ledger" ]; then \
		docker exec -i xledger_postgres psql -U admin -d xledger -c "DROP TABLE IF EXISTS transactions, sub_categories, categories CASCADE;"; \
	elif [ "$(TABLE)" = "user" ]; then \
		docker exec -i xledger_postgres psql -U admin -d xledger -c "DROP TABLE IF EXISTS users CASCADE;"; \
		exit 1; \
	fi
	@echo "✅ Migration rolled back for $(TABLE)"

# 服务完整初始化（迁移 + 代码生成）
# 用法: make service-init TABLE=ledger
service-init:
	@if [ -z "$(TABLE)" ]; then \
		echo "❌ Error: TABLE parameter required. Usage: make service-init TABLE=ledger"; \
		exit 1; \
	fi
	@echo "🚀 Initializing $(TABLE) service..."
	@$(MAKE) migrate TABLE=$(TABLE)
	@$(MAKE) proto-rpc TABLE=$(TABLE)
	@$(MAKE) proto-api TABLE=$(TABLE)
	@$(MAKE) model-table TABLE=$(TABLE)
	@echo "✅ $(TABLE) service initialized"

# 启动 RPC 服务
# 用法: make run-rpc TABLE=ledger
run-rpc:
	@if [ -z "$(TABLE)" ]; then \
		echo "❌ Error: TABLE parameter required. Usage: make run-rpc TABLE=ledger"; \
		exit 1; \
	fi
	@cd service/$(TABLE)/rpc && go run $(TABLE).go

# 启动 API 服务
# 用法: make run-api TABLE=ledger
run-api:
	@if [ -z "$(TABLE)" ]; then \
		echo "❌ Error: TABLE parameter required. Usage: make run-api TABLE=ledger"; \
		exit 1; \
	fi
	@cd service/$(TABLE)/api && go run $(TABLE).go