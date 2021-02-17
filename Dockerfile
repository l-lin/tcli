FROM alpine:3.12

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT [ "/tcli" ]

WORKDIR /

COPY tcli .
