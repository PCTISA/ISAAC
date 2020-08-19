# Build Go App
FROM golang:alpine
RUN apk add --no-cache git gcc g++
ENV CGO_ENABLED=1
ENV GOOS=linux
WORKDIR /app
COPY . .
RUN go build -o bot .

# Build Docker Image
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=0 /app/bot .

CMD ./bot