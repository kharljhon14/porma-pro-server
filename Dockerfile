FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN apk add curl
RUN  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz


COPY . .
RUN go build -o main main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY start.sh .
COPY internal/db/migration ./migration

RUN chmod +x /app/start.sh
EXPOSE 8080

CMD [ "/app/main" ]
