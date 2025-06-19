FROM golang:1.24 AS builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -mod=vendor -o skill-rocks cmd/main.go

FROM alpine:3.18

ENV TZ=Europe/Moscow
RUN apk add --no-cache tzdata && \
    addgroup -g 101 usr && \
        adduser -H -u 101 -G usr -s /bin/sh -D usr

WORKDIR /opt/skill-rocks

COPY --from=builder /app/skill-rocks .
COPY --from=builder /app/pkg/db/migrations migrations


EXPOSE 8080

CMD ["./skill-rocks"]