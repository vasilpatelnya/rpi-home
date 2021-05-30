#BUILD STAGE
FROM golang:1.14.8-alpine3.12 as builder
WORKDIR /app/build
COPY . .
RUN go mod download
RUN go build -o rpihome ./cmd/rpihome

#PRODUCTION STAGE
FROM alpine:3.12.0
WORKDIR /app
COPY --from=builder /app/build/rpihome /app/rpihome
CMD mkdir backup
ENV ENVIRONMENT=docker
CMD ./rpihome