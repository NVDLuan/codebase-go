APP        ?= my-app
AWS_REGION ?= ap-southeast-1
ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text 2>/dev/null)
REGISTRY   ?= $(ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
TAG        ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
ENV        ?= dev

.PHONY: help test build push render plan-% apply-% fmt lint local cleanup

help: ## Liệt kê lệnh
	@grep -E '^[a-zA-Z_%-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  %-14s %s\n",$$1,$$2}'

test: ## Lint + kiểm tra kiến trúc + unit test của app
	cd app && go vet ./... && go test -race ./...

build: ## Build image (tag = git SHA)
	docker build -t $(APP):$(TAG) app

push: build ## Push image lên ECR
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(REGISTRY)
	docker tag $(APP):$(TAG) $(REGISTRY)/$(APP):$(TAG)
	docker push $(REGISTRY)/$(APP):$(TAG)

render: ## Render manifest cho ENV (vd: make render ENV=prod)
	kubectl kustomize deploy/overlays/$(ENV)

plan-%: ## terraform plan cho môi trường (vd: make plan-dev)
	terraform -chdir=infra/envs/$* init -input=false
	terraform -chdir=infra/envs/$* plan -out=tfplan

apply-%: ## terraform apply plan đã tạo (vd: make apply-dev)
	terraform -chdir=infra/envs/$* apply tfplan

fmt: ## Format Terraform
	terraform fmt -recursive infra

lint: ## Chạy toàn bộ pre-commit hook
	pre-commit run --all-files

local: ## Chạy app trên cluster kind local
	./scripts/local-dev.sh

cleanup: ## Dọn tài nguyên môi trường dev
	./scripts/cleanup.sh
