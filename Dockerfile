FROM alpine:latest
RUN apk add --update ca-certificates

ADD ./alertbot /usr/bin/alertbot

EXPOSE 8080

ENTRYPOINT ["/usr/bin/alertbot"]
