# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

ARG GOPROXY=https://goproxy.io,direct
ARG GOSUMDB=sum.golang.org
ENV GOPROXY=${GOPROXY} \
    GOSUMDB=${GOSUMDB}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/user-service ./cmd/app

# --- Runtime stage ---
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

COPY --from=builder /out/user-service /app/user-service
COPY ../.env /app/.env

ENV GIN_MODE=release \
    PORT=8081

EXPOSE 8081

RUN adduser -D -H appuser && chown -R appuser /app
USER appuser

ENTRYPOINT ["/app/user-service"]
