FROM golang:1.22.2-alpine3.19

WORKDIR /go/src/github.com/dmowcomber/dcc-go
COPY . /go/src/github.com/dmowcomber/dcc-go
RUN go install -mod=vendor

CMD "dcc-go"
