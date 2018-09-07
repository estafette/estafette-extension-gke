FROM google/cloud-sdk:215.0.0-alpine

LABEL maintainer="estafette.io" \
      description="The estafette-extension-gke component is an Estafette extension to run commands against Kubernetes Engine"

RUN apk add --update --no-cache gettext \
    && gcloud components install kubectl \
    && rm -rf /var/cache/apk/*

COPY estafette-extension-gke /
COPY templates /templates

ENTRYPOINT ["/estafette-extension-gke"]