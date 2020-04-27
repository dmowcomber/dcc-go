FROM golang:1.13.3-alpine3.10

WORKDIR /go/src/github.com/dmowcomber/dcc-go
COPY . /go/src/github.com/dmowcomber/dcc-go
RUN go install -mod=vendor

CMD "dcc-go"
