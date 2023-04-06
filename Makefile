# Go parameters
GOCMD=go
GORUN= $(GOCMD) run
GOTEST=$(GOCMD) test
GOCOVER= $(GOCMD) tool cover
GOCLEAN= $(GOCMD) clean

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
	$(GOTEST) -v $(TEST_FLAGS) ./...

cover:
	@echo "Calculating coverage..."
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCOVER) -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

clean:
	@echo "Cleaning up..."
	@$(GOCLEAN) $(CLEAN_FLAGS)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

.PHONY: all test cover clean