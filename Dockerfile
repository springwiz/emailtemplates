FROM golang:1.15.0-buster AS builder
ENV GO111MODULE=auto
WORKDIR /src
ADD . /src
RUN go build ./cmd/email
ENTRYPOINT ["./email","send"]