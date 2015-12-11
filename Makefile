# SHELL = /bin/bash

setup: initial build

initial:
	@echo "==== get gb"
	go get -u github.com/constabulary/gb/...
	@echo
	
	@echo "==== install npm"
	cd front && npm install && npm run setup
	@echo
	
build: clean
	@echo "==== build back"
	$(GOPATH)/bin/gb build
	@echo
	
	@echo "==== build front"
	cd front && npm run build
	cp -r front/website/ bin/front/
	@echo

clean:
	@if [ -d bin ]; then \
		echo "==== clean"; \
		rm -rd bin; \
		echo; \
	fi
