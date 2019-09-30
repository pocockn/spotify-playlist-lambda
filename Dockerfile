FROM vidsyhq/go-base:latest
LABEL maintainer="Nick Pocock"

ARG VERSION
LABEL version=$VERSION

ADD spotify-api /
ADD config /config

ENTRYPOINT ["/spotify-api"]