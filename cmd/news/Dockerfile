FROM alpine
RUN apk add --no-cache ca-certificates
COPY news /news
EXPOSE 443
ENTRYPOINT ["/news"]