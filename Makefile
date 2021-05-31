.PHONY: all
all:	public wu

mocks:
	mockgen -source internal/services/works_service.go -destination internal/mocks/works_service.go --package mocks
	mockgen -source internal/services/activities_service.go -destination internal/mocks/activities_service.go --package mocks
	mockgen -source internal/repositories/transaction_runner.go -destination internal/mocks/transaction_runner.go --package mocks
	mockgen -source internal/repositories/works_repository.go -destination internal/mocks/works_repository.go --package mocks
	mockgen -source internal/repositories/activities_repository.go -destination internal/mocks/activities_repository.go --package mocks
	mockgen -source internal/tools/file_uploader.go -destination internal/mocks/file_uploader.go --package mocks
	mockgen -source internal/tools/uuid_generator.go -destination internal/mocks/uuid_generator.go --package mocks

.PHONY: dev
dev:
	yarn dev

.PHONY: run
run: public
	air

.PHONY: test
test:
	go test -coverprofile=cover.out -v ./...
	go tool cover -html=cover.out -o cover.html
	yarn test

public: nuxt.config.js web
	yarn run generate

wu:
	go build cmd/wu/main.go
