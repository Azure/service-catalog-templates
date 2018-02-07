FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY service-catalog-templates /usr/local/bin/

RUN adduser -D svcatt
USER svcatt

ENTRYPOINT ["service-catalog-templates"]
