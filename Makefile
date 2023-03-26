GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: ## run format, vet, lint and test
	@make fmt vet lint test

.PHONY: dev
dev: ## install golangci-lint
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1

.PHONY: build
build: ## build cobra-test 
	@cd examples/cobra-test && go build . && mv cobra-test ../../bin

.PHONY: build-ksqlgrammar
build-ksqlgrammar: ## build ksqlgrammar
	@cd examples/ksqlgrammar && go build . && mv ksqlgrammar ../../bin

.PHONY: mockery
mockery: ## generate mockery http client
	# 
	# mockery --all --inpackage --keeptree
	# mockery --name Recognizer --srcpkg github.com/antlr/antlr4/runtime/Go/antlr
	# mockery --name Ksqldb 
	# mockery --name KsqldbFactory
	@mockery --name HTTPClient --keeptree --recursive
	# mockery --name NewClientWithOptionsFactory --keeptree --recursive
	# mockery --name NewClientFactory --keeptree --recursive
	# mockery --name TransportFactory --keeptree --recursive

.PHONY: test
test: ## run unit tests
	@$(GOTEST) -v ./... -short

.PHONY: test-sec
test-sec: ## run gosec
	@gosec ./...

.PHONY: test-cover
test-cover: ## run coverage tools
	@$(GOTEST) ./... -coverprofile=coverage/coverage.out
	@$(GOCOVER) -func=coverage/coverage.out 
	@$(GOCOVER) -html=coverage/coverage.out
	@golangci-lint run ./... --verbose --no-config --out-format checkstyle > coverage/golangci-lint.out

.PHONY: test-ci
test-ci: ## run linter
	@golangci-lint run ./... --verbose --no-config

.PHONY: vet
vet: ## run go vet on the source files
	@$(GO) vet ./...

.PHONY: doc
doc: ## generate godocs and start a local documentation webserver on port 8085
	@GO111MODULE=off godoc -notes=TODO -goroot=. -http=:8085 -index

.PHONY: lint
lint: ## run golangci lint
	@golangci-lint run

.PHONY: create-grammar
create-grammar: ## run antlr to create the ksql parser
	@java -jar `pwd`/antlr/antlr-4.12.0-complete.jar -Dlanguage=Go -o parser KSql.g4 

.PHONY: clean-compose
clean-compose: ## run docker-compose in cobra-test directory
	@cd examples/cobra-test && docker-compose down && docker-compose up -d && cd ..

.PHONY: fmt
fmt: ## run golang code formatter
	@$(GO) fmt ./...

.PHONY: changelog
changelog: ## create changelog
	@git-chglog --output CHANGELOG.md

.PHONY: shell
shell: ## execute ksqldb-cli in docker container
	@docker exec -it ksqldb-cli ksql http://ksqldb:8088
