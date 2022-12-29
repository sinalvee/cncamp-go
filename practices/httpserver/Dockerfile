FROM golang:1.19-alpine as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.14

ENV GIN_MODE=release \
    PORT=80

WORKDIR /app

COPY --from=builder /app/httpserver .

EXPOSE 80

ENTRYPOINT ["./httpserver"]