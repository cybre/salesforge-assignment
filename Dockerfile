FROM golang:1.21.7-alpine3.19 as builder
WORKDIR /app
COPY ./ ./
RUN go build -o server ./cmd/server

FROM alpine:3.19 AS prod
WORKDIR /app
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/server /app/
ENTRYPOINT ["/app/server"]