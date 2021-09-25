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
	mockgen -source internal/lib/storage_client.go -destination internal/mocks/storage_client.go --package mocks
	mockgen -source internal/lib/uuid_generator.go -destination internal/mocks/uuid_generator.go --package mocks

.PHONY: dev_front
dev_front:
	yarn dev

.PHONY: dev_back
dev_back: public
	air

.PHONY: test_unit_back
test_unit_back:
	go test -coverprofile=cover.out -v ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: test_unit_front
test_unit_front:
	yarn test

.PHONY: test_unit
test_unit: test_unit_back test_unit_front

test_dir=.test
.PHONY: test_api
test_api:
	mkdir -p $(test_dir)
	go run cmd/wu/main.go > /dev/null &
	@curl -f https://api.getpostman.com/collections/$(POSTMAN_COLLECTION_ID)?apikey=$(POSTMAN_API_KEY) > $(test_dir)/api.json
	@curl -f https://api.getpostman.com/environments/$(POSTMAN_ENVIRONMENT_ID)?apikey=$(POSTMAN_API_KEY) > $(test_dir)/env.json
	@curl -f -G -d 'key=$(PIXABAY_API_KEY)&id=$(PIXABAY_ID_THUMBNAIL)' https://pixabay.com/api/ | jq -r '.hits[0].largeImageURL' | xargs -I@ curl -o .test/test_thumb.jpg @
	@curl -f -G -d 'key=$(PIXABAY_API_KEY)&id=$(PIXABAY_ID_CONTENT)' https://pixabay.com/api/ | jq -r '.hits[0].largeImageURL' | xargs -I@ curl -o .test/test_content.jpg @
	newman run $(test_dir)/api.json -e $(test_dir)/env.json --working-dir $(test_dir);\
	result=$$?;\
	lsof -t -i:8000 | xargs kill -9
	rm -rf $(test_dir);\
	exit $$result

.PHONY: test
test: test_unit test_api

public: nuxt.config.js web
	yarn run generate

wu:
	go build cmd/wu/main.go

.PHONY: run
run:
	go run cmd/wu/main.go
