# Multi stages to reduce the image size since we only need the binary file to exec
# First stage: build stage
FROM golang:1.19-alpine3.16 AS firststage
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage:
FROM alpine:3.16
WORKDIR /app
COPY --from=firststage /app/main .
COPY --from=firststage /app/app.env .

EXPOSE 8080
CMD ["/app/main"]