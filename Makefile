# Set up variables
EXECUTABLE_NAME := mt-ddns-updater
SUBMODULES := ./pkg/dhtml ./pkg/dhtmlform ./pkg/dhtmlbs ./pkg/mbr ./internal/mtweb

#do the job
include internal/goapp/Makefile.inc.mk

### ADDITIONAL TARGETS
# build css and scss assets
.PHONY: assets
assets:
	@echo "Building assets"
	sass --stop-on-error --no-source-map --style=compressed --silence-deprecation "color-functions,global-builtin,if-function,import" web/assets/css/
	npx --quiet terser web/assets/js-src/main.js --mangle --compress --output web/assets/script.min.js

# add assets dependency to before-build step
before-build:: assets

# build and run
.PHONY: run
run: clear assets
	go build
	${EXECUTABLE_NAME} run
