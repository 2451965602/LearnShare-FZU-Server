
MODULE = LearnShare

DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output

# 启动必要的环境，比如 etcd、mysql
.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose-env.yml up -d

# 关闭必要的环境，但不清理 data（位于 docker/data 目录中）
.PHONY: env-down
env-down:
	@ cd ./docker && docker compose down


# 生成基于 Hertz 的脚手架
.PHONY: hz-%
hz-%:
	hz update -idl ${IDL_PATH}/$*.thrift

# 清除所有的构建产物
.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print

# 清除所有构建产物、compose 环境和它的数据
.PHONY: clean-all
clean-all: clean
	@echo "$(PREFIX) Checking if docker-compose services are running..."
	@docker compose -f ./docker/docker-compose.yml ps -q | grep '.' && docker compose -f ./docker/docker-compose.yml down || echo "$(PREFIX) No services are running."
	@echo "$(PREFIX) Removing docker data..."
	rm -rf ./docker/data

.PHONY: build
build:
	@docker build -t ghcr.io/2451965602/learnshare:latest .

.PHONY: run
run:
	@docker compose -f ./docker/docker-compose.yml pull
	@docker compose -f ./docker/docker-compose.yml up -d

.PHONY: stop
stop:
	@docker compose -f ./docker/docker-compose.yml down
