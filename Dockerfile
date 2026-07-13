FROM golang:1.26.5 AS build
WORKDIR /app

COPY go.mod go.sum ./
COPY engine/go.mod ./engine/

RUN go mod download

COPY . .
RUN go build -o /out/api ./apps/api

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=build /out/api /app/api

ENV DATABASE_PATH=/data/app.db

RUN mkdir -p /data
VOLUME ["/data"]

EXPOSE 8080

CMD ["/app/api"]