# Set up variables
EXECUTABLE_NAME := mt-ddns-updater
SUBMODULES := ./pkg/dhtml ./pkg/dhtmlform ./pkg/dhtmlbs ./pkg/mbr ./internal/mtweb

# add assets dependency to before-build step, clear and clean targets repeated to be in the beginning
before-build:: clear clean assets

# clean assets when cleaning workspace
clean:: clean-assets

# do the main job
include internal/goapp/Makefile.inc.mk

### PROJECT CUSTOM TARGETS
# build css and scss assets
.PHONY: assets
assets:
	@echo "Building assets"
	sass --stop-on-error --no-source-map --style=compressed --silence-deprecation "color-functions,global-builtin,if-function,import" web/assets/css/
	npx --quiet terser web/assets/js-src/main.js --mangle --compress --output web/assets/script.min.js

.PHONY: clean-assets
clean-assets:
	rm -f web/assets/css/*.css
	rm -f web/assets/*.min.js

# build and run
.PHONY: run
run: clear assets
	go build
	${EXECUTABLE_NAME} run
