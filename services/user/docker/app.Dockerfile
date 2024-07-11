# Builder
FROM golang:1.22.4 as builder
LABEL authors="arcorium"

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY .. .

RUN go mod tidy
RUN go mod download

RUN make build

# Run tester
FROM builder as test-runner

RUN go test ./...

# Runner
FROM alpine:latest as runner

WORKDIR /

COPY --from=builder /app /app

# Grpc Server
EXPOSE 8080
# Metric Server
EXPOSE 8081

ENTRYPOINT ["./app/build/server"]