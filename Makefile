.PHONY: help
help: ## Print make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: ent-install
ent-install: ## Install Ent code-generation module
	go get entgo.io/ent/cmd/ent

.PHONY: air-install
air-install: ## Install air
	go install github.com/air-verse/air@latest

.PHONY: ent-gen
ent-gen: ## Generate Ent code
	go generate ./ent

.PHONY: ent-new
ent-new: ## Create a new Ent entity (ie, make ent-new name=MyEntity)
	go run entgo.io/ent/cmd/ent new $(name)

.PHONY: admin
admin: ## Create a new admin user (ie, make admin email=myemail@web.com)
	go run cmd/admin/main.go --email=$(email)

.PHONY: run
run: ## Run the application
	clear
	go run cmd/web/main.go

.PHONY: watch
watch: ## Run the application and watch for changes with air to automatically rebuild
	clear
	air

.PHONY: test
test: ## Run all tests
	go test ./...

.PHONY: check-updates
check-updates: ## Check for direct dependency updates
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep "\["

.PHONY: docker-build
docker-build: ## Build the application
	nixpacks build . --name pagode
