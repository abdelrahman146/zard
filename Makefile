.PHONY: run build test test-shared generate-mocks

run:
	@echo "Running service: $(filter-out $@,$(MAKECMDGOALS))"
	go run services/$(filter-out $@,$(MAKECMDGOALS))/cmd/main.go

build:
	@echo "Building service: $(filter-out $@,$(MAKECMDGOALS))"
	go build -o bin/$(filter-out $@,$(MAKECMDGOALS)) services/$(filter-out $@,$(MAKECMDGOALS))/cmd/main.go

test:
	@echo "Testing service: $(filter-out $@,$(MAKECMDGOALS))"
	go test services/$(filter-out $@,$(MAKECMDGOALS))/...

test-shared:
	@echo "Running all tests"
	go test shared/...

SOURCE_DIRS := shared/validator shared/rpc shared/pubsub shared/cache shared/config
MOCKGEN := mockgen
generate-mocks:
	@echo "Generating mocks"
	@for dir in $(SOURCE_DIRS); do \
		source_file=$$dir/$$(basename $$dir).go; \
		mock_file=$$dir/$$(basename $$dir)_mock.go; \
		package_name=$$(basename $$dir); \
		$(MOCKGEN) -source=$$source_file -destination=$$mock_file -package=$$package_name; \
	done

