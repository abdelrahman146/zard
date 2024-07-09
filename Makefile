.PHONY: run build

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