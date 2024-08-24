LOCAL_BIN:=$(CURDIR)/bin
PROTOC:=$(LOCAL_BIN)/protoc-27.3-linux-x86_64/bin/protoc
CONTRACTS_DIR:=$(CURDIR)/contracts

install:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

generate:
	make generate-sso-api

generate-sso-api:
	mkdir -p $(CONTRACTS_DIR)/gen/go
	$(PROTOC) -I $(CONTRACTS_DIR)/proto $(CONTRACTS_DIR)/proto/sso/sso.proto \
 	--go_out=$(CONTRACTS_DIR)/gen/go --go_opt=paths=source_relative \
 	--go-grpc_out=$(CONTRACTS_DIR)/gen/go --go-grpc_opt=paths=source_relative \
 	--plugin=protoc-gen-go=bin/protoc-gen-go \
 	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc