FROM alpine:3.7

WORKDIR /go/bin/

COPY ./app .
COPY ./config.yaml .

EXPOSE 8080

CMD ["./app"]
