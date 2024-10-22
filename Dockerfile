FROM golang:1.23.2 AS builder

COPY . /github.com/vbulash/auth/src/
WORKDIR /github.com/vbulash/auth/src/

RUN go mod download
RUN go build -o bin/auth_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/vbulash/auth/src/bin/auth_server .

# CMD ["./auth_server"]
