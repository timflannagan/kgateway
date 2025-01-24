FROM gcr.io/distroless/static:nonroot

ARG GOARCH=amd64

COPY sds-linux-$GOARCH /usr/local/bin/sds

USER 10101

ENTRYPOINT ["/usr/local/bin/sds"]
