# SHELL = /bin/bash

setup: initial build

initial:
	@echo "==== install npm"
	cd front && npm install && npm run setup
	@echo

build:
	@echo "==== build front"
	cd front && npm run build
	ln -s ../front/website bin/front
	@echo
