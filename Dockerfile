FROM golang:1.8-alpine
MAINTAINER Frederic Leger <leger.frederic@openmailbox.org>

RUN mkdir -p /go

# adds git to the image
RUN apk update && apk add git

COPY src/ /go/src/

# deps
RUN go-wrapper download gopkg.in/telegram-bot-api.v4

# install
RUN cd /go/src && go-wrapper install github.com/fredleger/golang/parrot
RUN cd /go/src && go-wrapper install github.com/fredleger/golang/parrotbot

WORKDIR /go
CMD ["/go/bin/parrotbot"]
