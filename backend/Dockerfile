# Stage 1: Build (Keep this as is)
FROM golang:1.26-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Stage 2: Final compact image
FROM alpine:latest
WORKDIR /root/

# --- ADD THIS LINE ---
RUN apk add --no-cache tzdata
# ---------------------

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]