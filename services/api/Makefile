.PHONY: help run build test lint fmt tidy verify clean swagger smoke

APP      := yanmai-api
BIN_DIR  := bin
GO       := go

help: ## 显示可用命令
	@echo "研脉学术知识付费平台 · API 服务"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

run: ## 启动开发服务（自动执行迁移 + 种子数据）
	$(GO) run ./cmd/server

build: ## 编译并输出到 bin/
	$(GO) build -trimpath -o $(BIN_DIR)/$(APP).exe ./cmd/server

test: ## 运行全部测试并输出覆盖率
	$(GO) test ./... -cover

lint: ## 静态检查 + 格式校验
	$(GO) vet ./...
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then echo "以下文件未格式化:"; echo "$$unformatted"; exit 1; fi; \
	echo "格式检查通过"

fmt: ## 自动格式化
	gofmt -w .

tidy: ## 整理依赖
	$(GO) mod tidy

verify: lint test ## 一键校验（lint + test）

swagger: ## 生成 OpenAPI 文档到 docs/（前端 SDK 类型来源）
	$(GO) run github.com/swaggo/swag/cmd/swag@latest init \
		-g cmd/server/main.go -o docs --parseDependency --parseInternal --outputTypes json,yaml
	@echo "已生成 docs/swagger.json，接着可执行：pnpm --filter @yanmai/api-client gen:openapi"

smoke: ## 端到端冒烟测试（需服务已启动）
	bash ./scripts/smoke.sh

clean: ## 清理编译产物
	rm -rf $(BIN_DIR) coverage.out coverage.html
