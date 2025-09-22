FROM golang:1.24.0-alpine3.20 AS builder

RUN apk add --no-cache \
    git \
    build-base \
    ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o payment-service ./cmd/payment-service

FROM alpine:3.20

RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-dejavu \
    font-noto \
    font-noto-cjk \
    font-noto-emoji \
    tzdata

WORKDIR /app

COPY --from=builder /app/payment-service /app/payment-service

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_BIN=/usr/bin/chromium-browser

EXPOSE 8003

ENTRYPOINT ["/app/payment-service"]
