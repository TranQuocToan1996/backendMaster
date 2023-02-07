# Multi stages to reduce the image size since we only need the binary file to exec
# First stage: build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


# Run stage:
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration 

EXPOSE 8080
CMD ["/app/main"]
# /app/main /app/start.sh
ENTRYPOINT [ "/app/start.sh" ]