FROM alpine:latest as alpine
RUN apk add -U --no-cache ca-certificates

WORKDIR /

COPY main main

CMD ["/main"]