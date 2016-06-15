get-deps:
	./install.sh

tests:
	go test ./utils
	go test ./migration
	go test ./config
	go test ./models/users
	go test ./models/apps
	go test ./models/histories

.PHONY: tests
