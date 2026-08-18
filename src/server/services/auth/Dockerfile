FROM golang:1.26.6-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY pkg ./pkg
COPY services ./services

WORKDIR /workspace/services/auth
RUN go build -o /out/auth ./cmd/server

FROM alpine:3.22

RUN addgroup -S nova && adduser -S nova -G nova

COPY --from=builder /out/auth /usr/local/bin/auth

USER nova
EXPOSE 8081

CMD ["auth"]
