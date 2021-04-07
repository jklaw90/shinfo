FROM golang:1.15.2-alpine3.12 AS base-builder
WORKDIR /shinfo
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY cmd cmd/
COPY internal internal/
COPY pkg pkg/
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

FROM scratch as base-server
COPY --from=base-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base-builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO=/zoneinfo.zip


# Room Service
FROM base-builder as room-builder
RUN go build -o /go/bin/room-service cmd/room/main.go

FROM base-server as room-service
COPY --from=room-builder /go/bin/room-service /room-service
ENTRYPOINT ["/room-service"]


# Message Service
FROM base-builder as message-builder
RUN go build -o /go/bin/message-service cmd/message/main.go

FROM base-server as message-service
COPY --from=message-builder /go/bin/message-service /message-service
ENTRYPOINT ["/message-service"]

