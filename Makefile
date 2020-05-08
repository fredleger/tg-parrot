ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GOPATH := $(ROOT_DIR):$(HOME)/go

all: clean vendors build

deep-clean: clean
	rm -rf src/github.com/fredleger/CocoTelegramParrotBot/vendor

clean:
	rm -rf bin/*
	rm -rf pkg/*
	cd src/github.com/fredleger/CocoTelegramParrotBot && go clean

vendors:
	cd src/github.com/fredleger/CocoTelegramParrotBot && go mod vendor

build:
	cd src/github.com/fredleger/CocoTelegramParrotBot && \
		CGO_ENABLED=0 go build \
    		-i -v -a -installsuffix cgo -gcflags "all=-N -l" \
    		-o $(ROOT_DIR)/bin/CocoTelegramParrotBot .

install:
	cd src/github.com/fredleger/CocoTelegramParrotBot && go install

docker-image:
	docker build -t webofmars/tg-parrot:develop .
