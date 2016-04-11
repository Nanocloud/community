get-deps:
	@cd nanocloud && ./install.sh

tests:
	 go test ./nanocloud/utils
