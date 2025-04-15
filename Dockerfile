FROM golang:1.23 AS builder

WORKDIR /build

RUN mkdir ./data
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 🔧 Build a static binary with no glibc dependency
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authnode

FROM alpine

WORKDIR /app

COPY --from=builder /build/authnode /app/authnode
COPY secrets /app/secrets

CMD ["./authnode", "-db=./data"]