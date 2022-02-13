FROM golang:1.17-alpine as build

WORKDIR /app

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go build -o main main.go

# stage 2: Go 애플리케이션 바이너리 파일과 필요한 파일만 alpine 컨테이너에 복사
FROM alpine:3.12
RUN apk --update add ca-certificates

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 80

# 컨테이너가 시작할 때 파일 실행
CMD ["/app/main"]