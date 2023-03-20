FROM golang:1.19.7-alpine3.17

WORKDIR /go/src/github.com/dmowcomber/dcc-go
COPY . /go/src/github.com/dmowcomber/dcc-go
RUN go install -mod=vendor

CMD "dcc-go"
