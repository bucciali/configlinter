FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o configlinter ./cmd/

FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget

WORKDIR /app
COPY --from=builder /app/configlinter .

EXPOSE 8080 9090

ENTRYPOINT ["./configlinter"]
CMD ["--serve"]