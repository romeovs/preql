BIN = bin

GO_FILES = $(shell find . -name '*.go' | sed 's/^.\///' | grep -v '^examples')

# Print a message
m = printf "\033[34;1mcontent\033[0m %s\033[0m\n"

build: $(BIN)/preql

bin/preql: $(GO_FILES)
	@$m "Building preql..."
	@mkdir -p $(BIN)
	@go build -o $(BIN)/preql ./preql/main.go

examples: bin/preql
	@$m "Running examples"
	@$m "bin/preql ./examples/users"
	@bin/preql ./examples/users
