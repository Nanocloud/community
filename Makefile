get-deps:
	@cd nanocloud && ./install.sh

tests:
	 go test ./nanocloud/utils
	 go test ./nanocloud/migration
	 go test ./nanocloud/config

.PHONY: tests
