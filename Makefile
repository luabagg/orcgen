# Test
TEST_FLAGS=-race

# Coverage
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Clean
CLEAN_FLAGS=-testcache

all: test cover

test:
	@echo "Running tests..."
	go test -v $(TEST_FLAGS) ./...

cover:
	@echo "Calculating coverage..."
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

clean:
	@echo "Cleaning up..."
	go clean $(CLEAN_FLAGS)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

.PHONY: all test cover clean