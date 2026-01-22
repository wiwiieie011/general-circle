FROM golang:1.25 AS builder

WORKDIR /ticket-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.21

WORKDIR /ticket-service

RUN apk --no-cache add ca-certificates

COPY --from=builder /ticket-service/main .

EXPOSE 8082

CMD ["./main"]
# Этап сборки
FROM golang:1.25 AS builder

WORKDIR /ticket-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.21

WORKDIR /ticket-service

RUN apk --no-cache add ca-certificates

COPY --from=builder /ticket-service/main .

EXPOSE 8082

CMD ["./main"]
