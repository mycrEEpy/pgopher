FROM alpine:3.20
ARG BINARY_PATH
WORKDIR /opt/pgopher
COPY $BINARY_PATH pgopher
USER 65534:65534
ENTRYPOINT ["/opt/pgopher/pgopher"]
