# Builder stage
FROM golang:1.19.4-alpine3.17 AS builder
WORKDIR /build

## Install dependencies
RUN apk update && apk add --no-cache build-base git
COPY . . 
RUN go mod download

## Build application
RUN go build -tags static_all,musl -o main .


# Runner stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/main .
COPY --from=builder /build/app.env .

EXPOSE 5000
EXPOSE 5432

CMD ["./main"]  
