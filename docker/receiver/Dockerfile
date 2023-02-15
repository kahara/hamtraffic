FROM alpine:3.17.3

RUN apk add netcat-openbsd

# FIXME Why aren't we seeing the received packets in log
ENTRYPOINT ["/bin/sh", "-c", "nc -l -u -k 4739 | hexdump -C"]
