FROM golang:1.17.3-alpine3.15

WORKDIR /app
COPY . .

RUN apk update && \
  apk add git

EXPOSE 8080