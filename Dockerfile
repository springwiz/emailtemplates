FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /src
ADD . /src
RUN go build ./cmd/email

FROM scratch
COPY --from=builder /src/email /
ENTRYPOINT ["/email","send"]