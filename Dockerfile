FROM golang:1.16 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -o books-data .

FROM alpine:3.13

RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime
RUN echo "Asia/Bangkok" >  /etc/timezone

WORKDIR /usr/src/app

COPY --from=builder /src/books-data /usr/src/app/books-data
COPY --from=builder /src/.env /usr/src/app/.env

RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

EXPOSE 8080
CMD ["./books-data"]