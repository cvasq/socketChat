FROM golang:alpine AS builder

WORKDIR /gotty
RUN apk update && apk add git && apk add ca-certificates
ENV GOPATH="/gotty"
RUN go get github.com/yudai/gotty
WORKDIR /socketChat-client-cli
ADD . /socketChat-client-cli
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /gotty/bin/socketChat-client-cli


FROM alpine
EXPOSE 8080
RUN apk update && apk add git && apk add ca-certificates
WORKDIR /gotty/bin/
COPY --from=builder /gotty/bin/gotty .
COPY --from=builder /gotty/bin/socketChat-client-cli .
ENTRYPOINT ["/gotty/bin/gotty", "-w", "--reconnect", "/gotty/bin/socketChat-client-cli"]
