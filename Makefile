.PHONY: proto
proto:
	mkdir -p build
	go build -o build/protoc-gen-go-sql .
	DEBUG=true export PATH=$(CURDIR)/build/:$$PATH && buf generate