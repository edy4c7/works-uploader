.PHONY: all
all:	public wu

mocks:
	mockgen -source internal/services/works_service.go -destination internal/mocks/works_service.go --package mocks
	mockgen -source internal/services/activities_service.go -destination internal/mocks/activities_service.go --package mocks
	mockgen -source internal/services/users_service.go -destination internal/mocks/users_service.go --package mocks
	mockgen -source internal/repositories/transaction_runner.go -destination internal/mocks/transaction_runner.go --package mocks
	mockgen -source internal/repositories/works_repository.go -destination internal/mocks/works_repository.go --package mocks
	mockgen -source internal/repositories/activities_repository.go -destination internal/mocks/activities_repository.go --package mocks
	mockgen -source internal/repositories/users_repository.go -destination internal/mocks/users_repository.go --package mocks
	mockgen -source internal/lib/file_uploader.go -destination internal/mocks/file_uploader.go --package mocks
	mockgen -source internal/lib/uuid_generator.go -destination internal/mocks/uuid_generator.go --package mocks

.PHONY: dev_front
dev_front:
	yarn dev

.PHONY: dev_back
dev_back: public
	air

.PHONY: test_unit
test_unit:
	go test -coverprofile=cover.out -v ./...
	go tool cover -html=cover.out -o cover.html
	yarn test

.PHONY: test_api
test_api:
	mkdir -p .test
	curl https://api.getpostman.com/collections/$(POSTMAN_COLLECTION_ID)?apikey=$(POSTMAN_API_KEY) > .test/api.json
	curl https://api.getpostman.com/environments/$(POSTMAN_ENVIRONMENT_ID)?apikey=$(POSTMAN_API_KEY) > .test/env.json
	touch .test/test_thumb.jpg
	touch .test/test_content.jpg
	-@newman run .test/api.json -e .test/env.json --working-dir .test
	rm -rf .test

public: nuxt.config.js web
	yarn run generate

wu:
	go build cmd/wu/main.go

.PHONY: run
run:
	go run cmd/wu/main.go
