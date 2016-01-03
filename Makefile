# SHELL = /bin/bash

setup: initial build

initial:
	@echo "==== install npm"
	cd front && npm install && npm run setup
	@echo

build:
	@echo "==== build front"
	cd front && npm run build
	if [ ! -h "bin/front" ]; then mkdir -p bin && ln -s ../front/website/ bin/front; fi
	@echo
