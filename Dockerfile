FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/bot/
RUN go build 

FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/cmd/bot/bot .
CMD ["/app/bot"]
