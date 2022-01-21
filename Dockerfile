FROM golang:1.17-alpine3.15 AS builder

COPY . /github.com/morozvol/AuthService/cmd/authserver
WORKDIR /github.com/morozvol/AuthService/cmd/authserver

RUN  go build -o ./bin/main  github.com/morozvol/AuthService/cmd/authserver



FROM alpine:3.15
WORKDIR /root/

COPY --from=0 /github.com/morozvol/AuthService/cmd/authserver/bin/main .
COPY --from=0 /github.com/morozvol/AuthService/cmd/authserver/conf  conf/

EXPOSE 8050

CMD ["./main"]