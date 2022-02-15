VERSION := $(shell git describe --tags --dirty)
BUILD := $(shell date +%FT%T%z)

.PHONY: config clean dir all

all: clean bluelox web

dir:
	mkdir -p bin

clean:
	rm -rf bin

web: wasm dir
	cp -r web bin/web && \
	cd bin && \
	cp bluelox.wasm web/js && \
	cd web && \
	zip -r ../web.zip . && \
	cd .. && \
	rm -rf web

wasm: dir
	go generate ./... && \
	cd cmd/wasm && \
	GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w -X github.com/nanmu42/bluelox/version.Version=$(VERSION) -X github.com/nanmu42/bluelox/version.BuildDate=$(BUILD)" -o bluelox.wasm && \
	cp bluelox.wasm $(PWD)/bin/bluelox.wasm

bluelox: bluelox.bin

%.bin: dir
	go generate ./... && \
	cd cmd/$* && \
	go build -trimpath -ldflags "-s -w -X github.com/nanmu42/bluelox/version.Version=$(VERSION) -X github.com/nanmu42/bluelox/version.BuildDate=$(BUILD)" && \
	cp $* $(PWD)/bin/$*