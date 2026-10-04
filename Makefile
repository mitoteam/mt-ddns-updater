# Set up variables
EXECUTABLE_NAME := mt-ddns-updater
SUBMODULES := ./pkg/dhtml ./pkg/dhtmlform ./pkg/dhtmlbs ./pkg/mbr ./internal/mtweb

#do the job
include internal/goapp/Makefile.inc.mk

# additional targets
.PHONY: run
run:
	clear
	go build
	${EXECUTABLE_NAME} run
