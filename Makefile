# export GOPATH := /Users/frederic/Documents/_work/40- CocoTheParrotBot

all: clean install-dep build install

clean: 
	rm -rf bin/*

install-dep:
	go get gopkg.in/telegram-bot-api.v4

build:
	cd src/github.com/fredleger/golang/parrot && go build
	cd src/github.com/fredleger/golang/parrotbot && go build

install:
	cd src/github.com/fredleger/golang/parrot && go install
	cd src/github.com/fredleger/golang/parrotbot && go install

dockerize:
	docker build -t webofmars/tg-parrot:develop .