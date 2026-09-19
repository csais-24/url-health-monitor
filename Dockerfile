FROM golang:1.26.8-bookworm AS builder

WORKDIR /app

COPY . .

RUN go build -o /app/url-health-monitor ./cmd/url-health-monitor

FROM debian:bookworm-slim

RUN apt-get update \ 
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/url-health-monitor .


ENTRYPOINT ["./url-health-monitor"]