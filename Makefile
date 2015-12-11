# SHELL = /bin/bash

setup: initial build

initial:
	@echo "==== get gb"
	go get -u github.com/constabulary/gb/...
	@echo

build: clean
	@echo "==== build back"
	$(GOPATH)/bin/gb build
	@echo

clean:
	@if [ -d bin ]; then \
		echo "==== clean"; \
		rm -rd bin; \
		echo; \
	fi
