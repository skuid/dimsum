FROM alpine

VOLUME /usr/src/config
VOLUME /usr/src/passwords

EXPOSE 8080

RUN apk add -U ca-certificates

COPY dimsum /bin/

ENTRYPOINT ["/bin/dimsum"]

CMD ["--config", "/opt/config/dimsum/config.yaml"]