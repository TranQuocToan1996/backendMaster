containernameposgre=postgres12
rootuser=root
pw=mysecretpassword
dbname=simple_bank
dbpath=db/migration
network=networksimplebank

dockerpull:
	docker pull postgres:12-alpine

runpostgres:
	docker run --name $(containernameposgre) --network $(network) -e POSTGRES_USER=$(rootuser) -e POSTGRES_PASSWORD=$(pw) -p 5432:5432 -d postgres:12-alpine

stoppostgre:
	docker stop $(containernameposgre)

startpostgre:
	docker start $(containernameposgre)

psqlexec:
	docker exec -it $(containernameposgre) psql -U $(rootuser)

shellexec:
	docker exec -it $(containernameposgre) /bin/sh

createdb:
	docker exec -it $(containernameposgre) createdb --username=$(rootuser) --owner=$(rootuser) $(dbname) 

dropdb:
	docker exec -it $(containernameposgre) dropdb $(dbname) 

logs:
	docker logs $(containernameposgre)

createMigrate:
	migrate create -ext sql -dir $(dbpath) -seq init_schema

migrateup:
	@echo "update to newest"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose up

migrateupaws:
	migrate -path $(dbpath) -database "postgresql://root:7B6D9j9R8DIU3xlZT6fw@database-1.cvsputm32sxh.ap-northeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	@echo "update one more seq"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose up 1

migratedown:
	@echo "rollback all"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose down

migratedown1:
	@echo "rollback last"
	migrate -path $(dbpath) -database "postgresql://$(rootuser):$(pw)@localhost:5432/$(dbname)?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

tests:
	go test -v -cover -race -timeout 300s -count=1 ./...

testsfail:
	go test -v -cover -race -timeout 300s -count=1 ./... | grep FAIL

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/TranQuocToan1996/backendMaster/db/sqlc Store 

gitpush: tests
	git push

buildsimplebank:
	docker build -t simplebank:latest .

runsimplebank:
	docker run --name simplebank --network $(network) -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:mysecretpassword@postgres:5432/simple_bank?sslmode=disable" simplebank:latest

createnetwork:
	docker network create $(network)
	
genkeychacha:
	openssl rand -hex 64 | head -c 32

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

evans:
	evans --host localhost --port 9090 -r repl


proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto

	

.PHONY: dockerpull runpostgres stoppostgre startpostgre psqlexec shellexec createdb dropdb logs createMigrate migrateup sqlc tests server mockgen migrateup1 migratedown1 buildsimplebank runsimplebank db_docs db_schema proto evans




# ------------------------
GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=example
VERSION?=0.0.0
SERVICE_PORT?=3000
DOCKER_REGISTRY?= #if set it should finished by /
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

## Build:
build: ## Build your project and put the output binary in out/bin/
	mkdir -p out/bin
	GO111MODULE=on $(GOCMD) build -mod vendor -o out/bin/$(BINARY_NAME) .

clean: ## Remove build related file
	rm -fr ./bin
	rm -fr ./out
	rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

vendor: ## Copy of all packages needed to support builds and tests in the vendor directory
	$(GOCMD) mod vendor

watch: ## Run the code with cosmtrek/air to have automatic reload on changes
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	docker run -it --rm -w /go/src/$(PACKAGE_NAME) -v $(shell pwd):/go/src/$(PACKAGE_NAME) -p $(SERVICE_PORT):$(SERVICE_PORT) cosmtrek/air

## Test:
test: ## Run the tests of the project
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/AlekSi/gocov-xml
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	gocov convert profile.cov | gocov-xml > coverage.xml
endif

## Lint:
lint: lint-go lint-dockerfile lint-yaml ## Run all available linters

lint-dockerfile: ## Lint your Dockerfile
# If dockerfile is present we lint it.
ifeq ($(shell test -e ./Dockerfile && echo -n yes),yes)
	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--format checkstyle" || echo "" ))
	$(eval OUTPUT_FILE = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -i $(CONFIG_OPTION) hadolint/hadolint hadolint $(OUTPUT_OPTIONS) - < ./Dockerfile $(OUTPUT_FILE)
endif

lint-go: ## Use golintci-lint on your project
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)

lint-yaml: ## Use yamllint on the yaml file of your projects
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/thomaspoignant/yamllint-checkstyle
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | yamllint-checkstyle > yamllint-checkstyle.xml)
endif
	docker run --rm -it -v $(shell pwd):/data cytopia/yamllint -f parsable $(shell git ls-files '*.yml' '*.yaml') $(OUTPUT_OPTIONS)

## Docker:
docker-build: ## Use the dockerfile to build the container
	docker build --rm --tag $(BINARY_NAME) .

docker-release: ## Release the container with tag latest and version
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)