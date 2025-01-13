# syntax=docker/dockerfile:1
FROM golang:1.22.3-alpine3.18 as builder
WORKDIR /build
COPY . .
RUN <<EOF
  apk add --no-cache \
    ca-certificates \
    curl \
    git
EOF
RUN <<EOF
  CGO_ENABLED=0 \
    go build -v -o app
EOF

FROM alpine:3.21.2
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/app .
RUN chmod +x app
CMD ["./app"]
