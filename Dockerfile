FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ip-tracker main.go notion.go get_ipv4.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/ip-tracker .

ENTRYPOINT ["./ip-tracker"]

