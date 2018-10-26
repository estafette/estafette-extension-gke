FROM google/cloud-sdk:222.0.0-alpine

LABEL maintainer="estafette.io" \
      description="The estafette-extension-gke component is an Estafette extension to run commands against Kubernetes Engine"

RUN du -h / -d 1

RUN gcloud components install kubectl \
    gcloud components remove bq \
    gcloud components remove gsutil \
    && rm -rf /var/cache/apk/*

COPY estafette-extension-gke /
COPY templates /templates

RUN gcloud components list
RUN du -h / -d 1
RUN ls -latrh /
RUN ls -latrh /google-cloud-sdk

ENTRYPOINT ["/estafette-extension-gke"]