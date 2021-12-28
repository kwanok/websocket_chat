FROM golang:1.17.3-alpine3.15

WORKDIR /app
COPY . .

RUN apk update && \
  apk add git && \
  go get github.com/cespare/reflex

EXPOSE 8080
CMD ["reflex", "-c", "reflex.conf"]