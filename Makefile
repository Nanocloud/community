# SHELL = /bin/bash

setup: initial build

initial:
	@echo "==== get gb"
	go get -u github.com/constabulary/gb/...
	@echo

	@echo "==== install npm"
	cd front && npm install && npm run setup
	@echo

build: go npm

go:
	@echo "==== build back"
	$(GOPATH)/bin/gb build
	@echo

npm:
	@echo "==== build front"
	cd front && npm run build
	ln -s ../front/website bin/front
	@echo

