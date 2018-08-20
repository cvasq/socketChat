# Multi-stage build
FROM golang:alpine AS builder
WORKDIR /socketChat-server
ADD . /socketChat-server
RUN apk update && apk add git && apk add ca-certificates
RUN cd /socketChat-server && go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o socketChat-server

# Second stage, smaller image
FROM alpine
ENV HTTP_PORT=9001
RUN apk update && apk add git && apk add ca-certificates
WORKDIR /socketChat-server
COPY --from=builder /socketChat-server/socketChat-server .
ENTRYPOINT ./socketChat-server --port $HTTP_PORT
