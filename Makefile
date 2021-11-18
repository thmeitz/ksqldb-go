GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

.PHONY: fmt dev lint vet test test-cover build-cobra build-ksqlgrammar all

all:
	make fmt vet lint test

dev:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1


build:
	cd examples/cobra-test && go build . && mv cobra-test ../../bin

build-ksqlgrammar:  
	cd examples/ksqlgrammar && go build . && mv ksqlgrammar ../../bin

mockery:
	# 
	# mockery --all --inpackage --keeptree
	# mockery --name Recognizer --srcpkg github.com/antlr/antlr4/runtime/Go/antlr
	# mockery --name Ksqldb 
	# mockery --name KsqldbFactory
	# mockery --name HTTPClient --keeptree --recursive
	# mockery --name NewClientWithOptionsFactory --keeptree --recursive
	mockery --name NewClientFactory --keeptree --recursive
	# mockery --name TransportFactory --keeptree --recursive

test:
	$(GOTEST) -v ./... -short

test-sec:
	gosec ./...

test-cover:
	$(GOTEST) ./... -coverprofile=coverage/coverage.out
	$(GOCOVER) -func=coverage/coverage.out 
	$(GOCOVER) -html=coverage/coverage.out
	golangci-lint run ./... --verbose --no-config --out-format checkstyle > coverage/golangci-lint.out

test-ci:
	golangci-lint run ./... --verbose --no-config

vet:	## run go vet on the source files
	$(GO) vet ./...

doc:	## generate godocs and start a local documentation webserver on port 8085
	GO111MODULE=off godoc -notes=TODO -goroot=. -http=:8085 -index

lint:
	golangci-lint run

create-grammar:
	java -jar `pwd`/antlr/antlr-4.7.1-complete.jar -Dlanguage=Go -o parser KSql.g4 

clean-compose:	
	cd examples/cobra-test && docker-compose down && docker-compose up -d && cd ..

fmt: 
	$(GO) fmt ./...

changelog:
	git-chglog --output CHANGELOG.md