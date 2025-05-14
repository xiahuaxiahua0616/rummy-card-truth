# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
PROJ_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output

# ==============================================================================
# 定义其他需要的伪目标

.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -o $(OUTPUT_DIR)/ifonly-apiserver $(PROJ_ROOT_DIR)/cmd/ifonly-apiserver/main.go

.PHONY: format
format: # 格式化 Go 源码.
	@gofmt -s -w ./

.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@go mod tidy

.PHONY: clean
clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: debug
debug: build
	clear && $(OUTPUT_DIR)/ifonly-apiserver --mode debug

.PHONY: release
release: build
	clear && $(OUTPUT_DIR)/ifonly-apiserver --mode release