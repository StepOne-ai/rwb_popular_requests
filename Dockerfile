FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /popular-requests ./cmd/app

FROM scratch
COPY --from=builder /popular-requests /popular-requests
EXPOSE 8080
ENTRYPOINT ["/popular-requests"]
