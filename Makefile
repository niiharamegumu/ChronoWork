.PHONY: run run-prod build clean

SRC_DIR := .

run:
	 DATABASE_NAME=sqlite_dev.db go run $(SRC_DIR)/main.go

run-prod:
	 DATABASE_NAME=sqlite.db go run $(SRC_DIR)/main.go

build:
	go build -o chronowork main.go

clean:
	@rm -rf $(SRC_DIR)
