# ── Build stage ──
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go run cmd/packer/main.go
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ctf-tool main.go

# ── Runtime stage ──
FROM alpine:3.21
WORKDIR /app

RUN apk add --no-cache dumb-init nginx

COPY --from=build /out/ctf-tool /app/ctf-tool
COPY docker/entrypoint.sh /app/entrypoint.sh
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/install.html.template /app/install.html.template
RUN chmod 0755 /app/ctf-tool /app/entrypoint.sh

ENV PORT=8080 \
    QUIZ_BASE_URL=http://quiz.ktf.ninja \
    CTF_WEB_PORT=7681 \
    TERM=xterm-256color \
    COLORTERM=truecolor \
    LANG=C.UTF-8 \
    LC_ALL=C.UTF-8

EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/app/entrypoint.sh"]
