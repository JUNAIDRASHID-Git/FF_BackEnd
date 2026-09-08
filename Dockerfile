# Step 1: Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Step 2: Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/data ./data
COPY --from=builder /app/uploads ./uploads

EXPOSE 5050 10000

CMD ["./server"]
