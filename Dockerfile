FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/healthgate ./main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/healthgate /usr/local/bin/healthgate

EXPOSE 9115
ENTRYPOINT ["healthgate"]
CMD ["-config=/etc/healthgate/config.yaml", "-listen=:9115"]
