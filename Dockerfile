FROM alpine:latest as alpine
RUN apk add -U --no-cache ca-certificates

WORKDIR /

COPY beuss-server .

CMD ["/beuss-server"]