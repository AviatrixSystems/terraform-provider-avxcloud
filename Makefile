BIN = terraform-provider-avxcloud

build: $(BIN)

GO_FILES := $(shell find internal -name '*.go')

$(BIN): main.go $(GO_FILES)
	go build -o $(BIN)

.PHONY:
fmt:
	go fmt -w .

.PHONY:
generate::
	go generate ./...

.PHONY:
clean::
	rm -rf docs/
	rm -f $(BIN)
