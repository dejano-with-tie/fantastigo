.PHONY: build
build: build-fleet

.PHONY: build-fleet
build-fleet: openapi
	@./scripts/build.sh cmd/fleet/fleet.go ./main

.PHONY: openapi
openapi:
	@./scripts/openapi.sh fleet
