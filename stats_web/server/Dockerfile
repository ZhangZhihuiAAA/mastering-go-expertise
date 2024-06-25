FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

RUN mkdir $GOPATH/src/server
ADD ./stats.go $GOPATH/src/server
ADD ./handlers.go $GOPATH/src/server
WORKDIR $GOPATH/src/server
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod init
RUN go mod tidy
RUN mkdir /pro
RUN go build -o /pro/server stats.go handlers.go

FROM alpine:latest

WORKDIR /pro
COPY --from=builder /pro/server .
EXPOSE 1234
CMD ["/pro/server"]