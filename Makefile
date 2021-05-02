.PHONY: all
all:	
	cd web \
	&& yarn \
	&& yarn generate
	cd functions \
	&& yarn \
	&& yarn build

emulator: all
	firebase emulators:start
	&& yarn run build

.PHONY: dev
dev:
	cd web && yarn dev

.PHONY: test_web
test_web:
	cd web && yarn test

.PHONY: test_go
test_go:
	go test -coverprofile=cover.out -v ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: test
test: test_go test_web

public: nuxt.config.js web
	yarn run generate

mocks:
	mockgen -source internal/services/works_service.go -destination internal/mocks/works_service.go --package mocks
	mockgen -source internal/services/activities_service.go -destination internal/mocks/activities_service.go --package mocks
	mockgen -source internal/repositories/transaction_runner.go -destination internal/mocks/transaction_runner.go --package mocks
	mockgen -source internal/repositories/works_repository.go -destination internal/mocks/works_repository.go --package mocks
	mockgen -source internal/repositories/activities_repository.go -destination internal/mocks/activities_repository.go --package mocks
	mockgen -source internal/tools/file_uploader.go -destination internal/mocks/file_uploader.go --package mocks
	mockgen -source internal/tools/uuid_generator.go -destination internal/mocks/uuid_generator.go --package mocks

.PHONY: run
run: public
	go run cmd/dsw/main.go