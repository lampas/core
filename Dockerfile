FROM golang:1.22.5-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o pb main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/pb .
COPY seeds_data /app/seeds_data
EXPOSE 8090
CMD ["/app/pb", "serve", "--http=0.0.0.0:8090"]
