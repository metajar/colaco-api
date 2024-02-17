SPEC_FILE := ./internal/api/v1/api.yml
CONFIG_FILE := ./internal/api/v1/config.yaml
OUTPUT := ./internal/api/v1/api.gen.go
PKG := v1
OAPI_CODEGEN := $(shell go env GOPATH)/bin/oapi-codegen

gen:
	$(OAPI_CODEGEN) -config  $(CONFIG_FILE) -o $(OUTPUT) $(SPEC_FILE)

clean:
	rm $(OUTPUT)

docker:
	docker build -t cola .

run:
	docker run -p 8080:8080 cola

cli:
	go build -o cola-client ./cmd/client


.PHONY: generate
