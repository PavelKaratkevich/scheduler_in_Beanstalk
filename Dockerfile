# Builder stage
FROM golang:alpine AS builder
WORKDIR /build

## Install dependencies
RUN apk update && apk add --no-cache build-base git
COPY . ./
RUN go mod download

## Build application
RUN go build -tags static_all,musl -o main ./cmd/main.go


# Runner stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/main .

EXPOSE 5000

CMD ["./main"]  
