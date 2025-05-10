BINARY   := actionai
CMD_PATH := ./cmd/cli

.PHONY: all build run clean

all: build

build:
	go build -o bin/$(BINARY) $(CMD_PATH)
	cp ./actions.json ./bin/

run:
	go run $(CMD_PATH)

clean:
	rm -rf bin/


