.PHONY: help run build test lint fmt tidy verify clean swagger sdk smoke

APP      := yanmai-api
BIN_DIR  := bin
GO       := go
# golangci-lint 可执行文件。留空则自动探测（PATH → $(go env GOPATH)/bin），
# 需要显式指定时：make lint GOLANGCI=/path/to/golangci-lint
GOLANGCI ?=
# 固定 swag 版本：契约校验依赖生成结果逐字节稳定，本地与 CI 必须一致
SWAG_VERSION ?= v1.16.6

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

lint: ## 静态检查（golangci-lint；未加入 PATH 时自动回落到 $(go env GOPATH)/bin）
	@GCL="$(GOLANGCI)"; \
	if [ -z "$$GCL" ]; then \
		if command -v golangci-lint >/dev/null 2>&1; then GCL=golangci-lint; \
		else \
			for cand in "$$(go env GOPATH)/bin/golangci-lint" "$$(go env GOPATH)/bin/golangci-lint.exe"; do \
				if [ -x "$$cand" ]; then GCL="$$cand"; break; fi; \
			done; \
		fi; \
	fi; \
	if [ -z "$$GCL" ]; then \
		echo "未找到 golangci-lint。请安装，或用 make lint GOLANGCI=/path/to/golangci-lint 指定。" >&2; \
		exit 1; \
	fi; \
	echo "使用 $$GCL"; \
	exec "$$GCL" run ./...

fmt: ## 自动格式化
	$(GO) fmt ./...
	@command -v gofmt >/dev/null && gofmt -l . | grep . && echo "以上文件已自动格式化" || true

tidy: ## 整理依赖
	$(GO) mod tidy

verify: lint test ## 一键校验（lint + test）

swagger: ## 生成 OpenAPI 文档到 docs/（前端 SDK 类型来源）
	$(GO) run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION) init \
		-g cmd/server/main.go -o docs --parseDependency --parseInternal --outputTypes json,yaml

sdk: swagger ## 重新生成 OpenAPI 文档并同步前端类型定义
	cd "$(CURDIR)/../.." && pnpm --filter @yanmai/api-client gen:openapi
	@echo "前端类型已同步；请将 packages/api-client/src/generated/schema.d.ts 一并提交"

smoke: ## 端到端冒烟测试（需服务已启动）
	bash ./scripts/smoke.sh

clean: ## 清理编译产物
	rm -rf $(BIN_DIR) coverage.out coverage.html
