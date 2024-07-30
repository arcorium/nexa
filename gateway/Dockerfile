# Builder
FROM golang:1.22.4 AS builder
LABEL authors="arcorium"

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY docker .

RUN go mod tidy
RUN go mod download

RUN go build -o build/server "./cmd/server/"

# Runner
FROM alpine:latest AS runner

COPY --from=builder /app/build/* /app/

WORKDIR /app

# Gateway server
EXPOSE 8080

ENTRYPOINT ["./server"]