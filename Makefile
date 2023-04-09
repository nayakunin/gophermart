# Define variables
SCHEMA_FILE := ./api/schema.yaml
GEN_FILE := ./internal/generated/api.go
PACKAGE_NAME := api

ACCRUAL_SCHEMA_FILE := ./api/services/accrual.yaml
ACCRUAL_GEN_FILE := ./internal/services/accrual/generated/api.go
ACCRUAL_PACKAGE_NAME := api

LOCAL_BIN:=$(CURDIR)/bin
GOAPI_CODEGEN_BIN:=$(LOCAL_BIN)/goapi-gen

# install oapi-codegen tool
.PHONY: .install-goapi-gen
.install-opapi-codegen:
	$(info #Install goapi-gen tool)
	tmp=$$(mktemp -d) && cd $$tmp && pwd && \
		go mod init temp && \
		(go get -d github.com/discord-gophers/goapi-gen@latest) && \
		go build -o $(GOAPI_CODEGEN_BIN) github.com/discord-gophers/goapi-gen && \
		rm -rf $$tmp

# generate code from schema
.PHONY: all
generate: .install-opapi-codegen
	$(GOAPI_CODEGEN_BIN) -generate types -package $(ACCRUAL_PACKAGE_NAME) -out $(ACCRUAL_GEN_FILE) $(ACCRUAL_SCHEMA_FILE)
	$(GOAPI_CODEGEN_BIN) -generate types,server -package $(PACKAGE_NAME) -out $(GEN_FILE) $(SCHEMA_FILE)
