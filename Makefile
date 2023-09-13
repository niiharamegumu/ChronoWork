.PHONY: default run clean

SRC_DIR := .

default: build run

run:
	 go run $(SRC_DIR)/main.go --env=dev

run-prod:
	 go run $(SRC_DIR)/main.go --env=prod

clean:
	@rm -rf $(SRC_DIR)
