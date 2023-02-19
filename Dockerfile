# Build stage
FROM golang:1.20.1-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o exchangeApp main.go

RUN chmod +x /app/exchangeApp

# Run stage
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/exchangeApp .

COPY db/migration db/migration

COPY app.env .

CMD [ "/app/exchangeApp" ]
