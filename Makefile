.PHONY: default run clean

SRC_DIR := .

default: build run

run:
	 go run $(SRC_DIR)/main.go

clean:
	@rm -rf $(SRC_DIR)
