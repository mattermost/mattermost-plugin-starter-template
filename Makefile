GO ?= $(shell command -v go 2> /dev/null)
NPM ?= $(shell command -v npm 2> /dev/null)
CURL ?= $(shell command -v curl 2> /dev/null)
MANIFEST_FILE ?= plugin.json
GOPATH ?= $(shell go env GOPATH)
GO_TEST_FLAGS ?= -race
GO_BUILD_FLAGS ?=
MM_UTILITIES_DIR ?= ../mattermost-utilities

export GO111MODULE=on

# You can include assets this directory into the bundle. This can be e.g. used to include profile pictures.
ASSETS_DIR ?= assets
PWD := $(shell pwd)
OS := $(shell uname 2> /dev/null)

# Verify environment, and define PLUGIN_ID, PLUGIN_VERSION, HAS_SERVER and HAS_WEBAPP as needed.
include build/setup.mk

BUNDLE_NAME ?= $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz

# Include custom makefile, if present
ifneq ($(wildcard build/custom.mk),)
	include build/custom.mk
endif

## Checks the code style, tests, builds and bundles the plugin.
all: check-style test dist

## Propagates plugin manifest information into the server/ and webapp/ folders as required.
.PHONY: apply
apply:
	./build/bin/manifest apply

## Runs govet and gofmt against all packages.
.PHONY: check-style
check-style: webapp/.npminstall gofmt govet golint
	@echo Checking for style guide compliance

ifneq ($(HAS_WEBAPP),)
	cd webapp && npm run lint
endif

## Runs gofmt against all packages.
.PHONY: gofmt
gofmt:
ifneq ($(HAS_SERVER),)
	@echo Running gofmt
	@for package in $$(go list ./...); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "Gofmt failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo Gofmt success
endif

## Runs govet against all packages.
.PHONY: govet
govet:
ifneq ($(HAS_SERVER),)
	@echo Running govet
	@# Workaround because you can't install binaries without adding them to go.mod
	env GO111MODULE=off $(GO) get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	$(GO) vet ./...
	$(GO) vet -vettool=$(GOPATH)/bin/shadow ./...
	@echo Govet success
endif

## Runs golint against all packages.
.PHONY: golint
golint:
	@echo Running lint
	env GO111MODULE=off $(GO) get golang.org/x/lint/golint
	$(GOPATH)/bin/golint -set_exit_status ./...
	@echo lint success

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint on codebase
# https://stackoverflow.com/a/677212/1027058 (check if a command exists or not)
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi; \

	@echo Running golangci-lint
	golangci-lint run ./...

## Generate mocks.
mocks:
ifneq ($(HAS_SERVER),)
	go install github.com/golang/mock/mockgen
endif

## Builds the server, if it exists, including support for multiple architectures.
.PHONY: server
server:
ifneq ($(HAS_SERVER),)
	mkdir -p server/dist;
	cd server && env GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-linux-amd64;
	cd server && env GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-darwin-amd64;
	cd server && env GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-windows-amd64.exe;
endif

## Ensures NPM dependencies are installed without having to run this all the time.
webapp/.npminstall:
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) install
	touch $@
endif

## Builds the webapp, if it exists.
.PHONY: webapp
webapp: webapp/.npminstall
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) run build;
endif

## Generates a tar bundle of the plugin for install.
.PHONY: bundle
bundle:
	rm -rf dist/
	mkdir -p dist/$(PLUGIN_ID)
	cp $(MANIFEST_FILE) dist/$(PLUGIN_ID)/
ifneq ($(wildcard $(ASSETS_DIR)/.),)
	cp -r $(ASSETS_DIR) dist/$(PLUGIN_ID)/
endif
ifneq ($(HAS_PUBLIC),)
	cp -r public/ dist/$(PLUGIN_ID)/
endif
ifneq ($(HAS_SERVER),)
	mkdir -p dist/$(PLUGIN_ID)/server/dist;
	cp -r server/dist/* dist/$(PLUGIN_ID)/server/dist/;
endif
ifneq ($(HAS_WEBAPP),)
	mkdir -p dist/$(PLUGIN_ID)/webapp/dist;
	cp -r webapp/dist/* dist/$(PLUGIN_ID)/webapp/dist/;
endif
	cd dist && tar -cvzf $(BUNDLE_NAME) $(PLUGIN_ID)

	@echo plugin built at: dist/$(BUNDLE_NAME)

## Builds and bundles the plugin.
.PHONY: dist
dist:	apply server webapp bundle

## Installs the plugin to a (development) server.
## It uses the API if appropriate environment variables are defined,
## and otherwise falls back to trying to copy the plugin to a sibling mattermost-server directory.
.PHONY: deploy
deploy: dist
	./build/bin/pluginctl deploy $(PLUGIN_ID) dist/$(BUNDLE_NAME)

## Runs any lints and unit tests defined for the server and webapp, if they exist.
.PHONY: test
test: webapp/.npminstall
ifneq ($(HAS_SERVER),)
	$(GO) test -v $(GO_TEST_FLAGS) ./server/...
endif
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) run fix && $(NPM) run test;
endif

## Creates a coverage report for the server code.
.PHONY: coverage
coverage: webapp/.npminstall
ifneq ($(HAS_SERVER),)
	$(GO) test $(GO_TEST_FLAGS) -coverprofile=server/coverage.txt ./server/...
	$(GO) tool cover -html=server/coverage.txt
endif

## Extract strings for translation from the source code.
.PHONY: i18n-extract
i18n-extract:
ifneq ($(HAS_WEBAPP),)
ifeq ($(HAS_MM_UTILITIES),)
	@echo "You must clone github.com/mattermost/mattermost-utilities repo in .. to use this command"
else
	cd $(MM_UTILITIES_DIR) && npm install && npm run babel && node mmjstool/build/index.js i18n extract-webapp --webapp-dir $(PWD)/webapp
endif
endif

## Clean removes all build artifacts.
.PHONY: clean
clean:
	rm -fr dist/
ifneq ($(HAS_SERVER),)
	rm -fr server/dist
endif
ifneq ($(HAS_WEBAPP),)
	rm -fr webapp/.npminstall
	rm -fr webapp/dist
	rm -fr webapp/node_modules
endif
	rm -fr build/bin/

# Help documentation à la https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@cat Makefile | grep -v '\.PHONY' |  grep -v '\help:' | grep -B1 -E '^[a-zA-Z0-9_.-]+:.*' | sed -e "s/:.*//" | sed -e "s/^## //" |  grep -v '\-\-' | sed '1!G;h;$$!d' | awk 'NR%2{printf "\033[36m%-30s\033[0m",$$0;next;}1' | sort


##
## The following are commands that can make a developer's life easier
##

# sd is an easier-to-type alias for server-debug
.PHONY: sd
sd: server-debug

# server-debug builds and deploys a debug version of the plugin for your architecture.
# Then resets the plugin to pick up the changes.
.PHONY: server-debug
server-debug: apply server-debug-build reset

.PHONY: server-debug-build
server-debug-build:
	mkdir -p server/dist

ifeq ($(OS),Darwin)
	cd server && env GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -gcflags "all=-N -l" -o dist/plugin-darwin-amd64;
else ifeq ($(OS),Linux)
	cd server && env GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -gcflags "all=-N -l" -o dist/plugin-linux-amd64;
else ifeq ($(OS),Windows_NT)
	cd server && env GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -gcflags "all=-N -l" -o dist/plugin-windows-amd64.exe;
else
	$(error make debug depends on uname to return your OS. If it does not return 'Darwin' (meaning OSX), 'Linux', or 'Windows_NT' (all recent versions of Windows), you will need to edit the Makefile for your own OS.)
endif

	rm -rf dist/
	mkdir -p dist/$(PLUGIN_ID)/server/dist
	cp $(MANIFEST_FILE) dist/$(PLUGIN_ID)/
	cp -r server/dist/* dist/$(PLUGIN_ID)/server/dist/
	mkdir -p ../mattermost-server/plugins
	cp -r dist/* ../mattermost-server/plugins/

# wd is an easier-to-type alias for webapp-debug
.PHONY: wd
wd: webapp-debug

# webapp-debug builds and deploys the plugin's webapp in watch mode with source-maps.
# Webpack will run `make reset` after detecting and compiling changes.
.PHONY: webapp-debug
webapp-debug:
ifneq ($(HAS_WEBAPP),)
# link the webapp directory
	rm -rf ../mattermost-server/plugins/$(PLUGIN_ID)/webapp
	mkdir -p ../mattermost-server/plugins/$(PLUGIN_ID)/webapp
	ln -nfs $(PWD)/webapp/dist ../mattermost-server/plugins/$(PLUGIN_ID)/webapp/dist
# compile webapp with development options, and continuously watch for changes
	cd webapp && $(NPM) run debug:watch
endif

# wdd is an easier-to-type alias for webapp-debug-deploy
.PHONY: wdd
wdd: webapp-debug-deploy

# webapp-debug-deploy builds and deploys the plugin's webapp with source-maps.
# This is a one-time deploy similar to `make deploy`
.PHONY: webapp-debug-deploy
webapp-debug-deploy:
ifneq ($(HAS_WEBAPP),)
# link the webapp directory
	rm -rf ../mattermost-server/plugins/$(PLUGIN_ID)/webapp
	mkdir -p ../mattermost-server/plugins/$(PLUGIN_ID)/webapp
	ln -nfs $(PWD)/webapp/dist ../mattermost-server/plugins/$(PLUGIN_ID)/webapp/dist
# compile webapp with development options
	cd webapp && $(NPM) run debug
# re-enable plugin via API
	./build/bin/pluginctl reset $(PLUGIN_ID)
endif

# Reset the plugin
# If we were debugging, we have to unattach the delve process or else we can't disable the plugin.
# NOTE: we are assuming the dlv was listening on port 2346, as in the debug-plugin.sh script.
.PHONY: reset
reset:
	@DELVE_PID=$(shell ps aux | grep "dlv attach.*2346" | grep -v "grep" | awk -F " " '{print $$2}') && \
	if [ "$$DELVE_PID" -gt 0 ] > /dev/null 2>&1 ; then \
		echo "Located existing delve process running with PID: $$DELVE_PID. Killing." ; \
		kill -9 $$DELVE_PID ; \
	fi

	./build/bin/pluginctl reset $(PLUGIN_ID)

.PHONY: debug-plugin
debug-plugin:
	$(eval PLUGIN_PID := $(shell ps aux | grep "plugins/${PLUGIN_ID}" | grep -v "grep" | awk -F " " '{print $$2}'))
	$(eval NUM_PID := $(shell echo -n ${PLUGIN_PID} | wc -w))

	@if [ ${NUM_PID} -gt 2 ]; then \
		echo "** There is more than 1 plugin process running. Run 'make kill-plugin' to get rid of them."; \
		echo "   Then run 'make reset' to start the plugin process again, and 'make debug-plugin' attach the dlv process."; \
		exit 1; \
	fi

	@if [ -z ${PLUGIN_PID} ]; then \
		echo "Could not find plugin PID; the plugin is not running. Exiting."; \
		exit 1; \
	fi

	@echo "Located Plugin running with PID: ${PLUGIN_PID}"
	dlv attach ${PLUGIN_PID} --listen :2346 --headless=true --api-version=2 --accept-multiclient &

.PHONY: kill-plugin
kill-plugin:
# If we were debugging, we have to unattach the delve process or else we can't disable the plugin.
# NOTE: we are assuming the dlv was listening on port 2346, as in the debug-plugin.sh script.
	$(eval DELVE_PID := $(shell ps aux | grep "dlv attach.*2346" | grep -v "grep" | awk -F " " '{print $$2}'))

	@if [ -n "${DELVE_PID}" ]; then \
		echo "Located existing delve process running with PID: ${DELVE_PID}. Killing."; \
		kill -9 ${DELVE_PID}; \
	fi

	$(eval PLUGIN_PID := $(shell ps aux | grep "plugins/${PLUGIN_ID}" | grep -v "grep" | awk -F " " '{print $$2}'))

	@for PID in ${PLUGIN_PID}; do \
		echo "Killing plugin pid $$PID"; \
		kill -9 $$PID; \
	done; \
