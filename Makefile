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
	goctl model pg datasource -url="$(DB_URL)" -table="$(TABLE)" -dir="./service/$(TABLE)/model" -cache=true

# 构建ledger服务相关的model
model-all:
	@for table in $(TABLES); do \
		echo "正在生成表: $$table"; \
		mkdir -p ./service/$$table/model && \
		goctl model pg datasource -url="$(DB_URL)" -table="$$table" -dir="./service/$$table/model" -cache=true || true; \
	done

# 构建单个PRC代码
proto-rpc:
	@if [ -z "$(TABLE)" ] ; then \
		exit 1; \
	fi
	goctl rpc protoc service/$(TABLE)/rpc/$(TABLE).proto --go_out=service/$(TABLE)/rpc --go-grpc_out=service/$(TABLE)/rpc --zrpc_out=service/$(TABLE)/rpc

# 构建单个API代码
proto-api:
	@if [ -z "$(TABLE)" ] ; then \
		exit 1; \
	fi
	goctl api go -api service/$(TABLE)/api/$(TABLE).api -dir service/$(TABLE)/api