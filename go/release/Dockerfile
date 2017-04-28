
FROM alpine:3.3

RUN apk --update upgrade && apk add ca-certificates && update-ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
COPY app /app

EXPOSE 8080

ENTRYPOINT "/app/iot.ericslien.com"
