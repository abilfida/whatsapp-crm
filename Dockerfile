FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api

FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/.env.example /app/.env.example
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/api"]
