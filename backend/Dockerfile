FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
#COPY go.mod .

RUN <<EOF
go mod download
go mod tidy 
go build -o web-chat-bot
EOF

#FROM scratch
FROM debian:12-slim
# ------------------------------------
# Install Socat to use MCP Toolkit
# ------------------------------------
RUN <<EOF
apt-get update
apt-get install -y socat
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF


WORKDIR /app
COPY --from=builder /app/web-chat-bot .
COPY docs /app/docs

CMD ["./web-chat-bot"]