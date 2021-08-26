#BUILD STAGE
FROM golang:alpine as builder
RUN apk add --update gcc musl-dev
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
RUN mkdir db
ENV ENVIRONMENT=docker
CMD ./rpihome