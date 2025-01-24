FROM scratch

ARG GOARCH

COPY sds-linux-${GOARCH} /usr/local/bin/sds

USER 10101

ENTRYPOINT ["/usr/local/bin/sds"]
