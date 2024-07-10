# Builder
FROM golang:1.22.4 AS builder
LABEL authors="arcorium"

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .

RUN go mod tidy
RUN go mod download

RUN go build -o build/seed "./cmd/seed/"

# Run tester
FROM builder AS test-runner

RUN go test ./...

# Runner
FROM alpine:latest AS runner

COPY --from=builder /app/build/* /app/

WORKDIR /app

ENTRYPOINT ["./seed"]