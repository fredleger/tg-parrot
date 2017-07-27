ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
GOPATH := $(HOME)/go:$(ROOT_DIR)

all: clean install-dep build install

debug:
	@echo "GOPATH: $(GOPATH)"

clean:
	rm -rf bin/*
	rm -rf pkg/*
	cd src/github.com/fredleger/CocoTelegramParrotBot/parrotlib && go clean
	cd src/github.com/fredleger/CocoTelegramParrotBot/cocobot && go clean

install-dep:
	cd src/github.com/fredleger/CocoTelegramParrotBot/parrotlib && go get -d -v
	cd src/github.com/fredleger/CocoTelegramParrotBot/cocobot && go get -d -v

build:
	cd src/github.com/fredleger/CocoTelegramParrotBot/parrotlib && go build
	cd src/github.com/fredleger/CocoTelegramParrotBot/cocobot && go build

install:
	cd src/github.com/fredleger/CocoTelegramParrotBot/parrotlib && go install
	cd src/github.com/fredleger/CocoTelegramParrotBot/cocobot && go install

docker-image:
	docker build -t webofmars/tg-parrot:develop -f $(ROOT_DIR)/Dockerfile .
