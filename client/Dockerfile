FROM golang:1.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY run-demo.sh ./run-demo.sh
RUN ls -l ./run-demo.sh
COPY . .

# 🔧 Build a static binary with no glibc dependency
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-client

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /build/auth-client /app/auth-client
COPY --from=builder /build/run-demo.sh /app/run-demo.sh

RUN chmod +x /app/run-demo.sh

RUN ls -lh /app
CMD ["sh", "-c", "/app/run-demo.sh && sleep infinity"]