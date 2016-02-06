FROM alpine:3.3
MAINTAINER 0rax <jp@roemer.im>

RUN apk add -U ca-certificates

ADD go-redirect /app/go-redirect
COPY docker/nsswitch.conf /etc/nsswitch.conf

WORKDIR /app
ENTRYPOINT [ "/app/go-redirect" ]
