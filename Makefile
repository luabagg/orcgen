# Go parameters
GOCMD=go
GORUN= $(GOCMD) run
GOTEST=$(GOCMD) test
GOCOVER= $(GOCMD) tool cover
GOCLEAN= $(GOCMD) clean

# Test
TEST_FLAGS=-race

# Coverage
COVERAGE_FILE=testdata/coverage.out
COVERAGE_HTML=testdata/coverage.html

# Examples
EXAMPLE=webpage

all: test cover

test:
	@echo "Running tests..."
	$(GOTEST) -v $(TEST_FLAGS) ./...

cover:
	@echo "Calculating coverage..."
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCOVER) -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

example:
	@echo "Running example..."
	$(GORUN) examples/convert-$(EXAMPLE)/main.go

clean:
	@echo "Cleaning up..."
	@$(GOCLEAN) -testcache
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

.PHONY: all test cover clean