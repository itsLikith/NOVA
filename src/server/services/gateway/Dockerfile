FROM golang:1.26.6-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY pkg ./pkg
COPY services ./services

WORKDIR /workspace/services/gateway
RUN go build -o /out/gateway ./cmd/server

FROM alpine:3.22

RUN addgroup -S nova && adduser -S nova -G nova

COPY --from=builder /out/gateway /usr/local/bin/gateway

USER nova
EXPOSE 8080

CMD ["gateway"]
