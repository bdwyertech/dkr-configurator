FROM golang:1.17-alpine as helper
WORKDIR /go/src/
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -ldflags="-s -w" -trimpath .

FROM gcr.io/distroless/static:nonroot-amd64

ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.title="bdwyertech/configurator" \
    org.opencontainers.image.description="Configuration rendering sidecar" \
    org.opencontainers.image.authors="Brian Dwyer <bdwyertech@github.com>" \
    org.opencontainers.image.url="https://hub.docker.com/r/bdwyertech/configurator" \
    org.opencontainers.image.source="https://github.com/bdwyertech/dkr-configurator.git" \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.created=$BUILD_DATE \
    org.label-schema.name="bdwyertech/configurator" \
    org.label-schema.description="Configuration rendering sidecar" \
    org.label-schema.url="https://hub.docker.com/r/bdwyertech/configurator" \
    org.label-schema.vcs-url="https://github.com/bdwyertech/dkr-configurator.git" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.build-date=$BUILD_DATE

COPY --from=helper /go/src/configurator /.
CMD ["/configurator"]
