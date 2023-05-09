.PHONY: build
build: install build-fleet

.PHONY: build-fleet
build-fleet: openapi
	@./scripts/build.sh cmd/fleet/fleet.go ./main

.PHONY: openapi
openapi: install
	@./scripts/openapi.sh fleet

.PHONY: install
install: download
	@./scripts/install.sh

.PHONY: download
download:
	@echo "Download go.mod dependencies"
	@go mod download