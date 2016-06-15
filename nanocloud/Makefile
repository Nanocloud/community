get-deps:
	./install.sh

tests:
	go test ./nanocloud/utils
	go test ./nanocloud/migration
	go test ./nanocloud/config
	go test ./nanocloud/models/users
	go test ./nanocloud/models/apps
	go test ./nanocloud/models/histories

.PHONY: tests
