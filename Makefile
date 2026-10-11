APP := oczko
BIN_DIR := bin

.PHONY: all build clean install test run
all: build

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP) ./src

run: build
	./$(BIN_DIR)/$(APP) $(FILE)

test:
	go test ./...

install:
	go build -o "$$(go env GOPATH)/bin/$(APP)" ./src

clean:
	rm -rf $(BIN_DIR)
