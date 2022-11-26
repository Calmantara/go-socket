# builder container
FROM golang:alpine AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARC=amd go build -o main ./...

# final container
FROM alpine:3.16 AS production
RUN mkdir /app
COPY --from=builder /app .
ENTRYPOINT ["./main"]
