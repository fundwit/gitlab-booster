FROM busybox:1.31.1-glibc

COPY gitlab-booster /
COPY buildInfo.json /
ENV GIN_MODE release
ENTRYPOINT ["/gitlab-booster"]