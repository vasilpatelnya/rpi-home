#BUILD STAGE
FROM golang:alpine as builder
RUN apk add --no-cache --initdb gcc
RUN export CC=gcc
WORKDIR /app/build
COPY . .
RUN go mod download
RUN go build -o rpihome ./cmd/rpihome

#PRODUCTION STAGE
FROM alpine:3.12.0
WORKDIR /app
COPY --from=builder /app/build/rpihome /app/rpihome
RUN mkdir backup
RUN mkdir config
ENV ENVIRONMENT=docker
CMD ./rpihome