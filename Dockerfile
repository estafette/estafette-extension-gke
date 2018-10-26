FROM google/cloud-sdk:222.0.0-alpine

LABEL maintainer="estafette.io" \
      description="The estafette-extension-gke component is an Estafette extension to deploy applications to a Kubernetes Engine cluster"

RUN du -hd1 /google-cloud-sdk/.[^.]* /google-cloud-sdk/*
RUN du -hd1 /google-cloud-sdk/bin/.[^.]* /google-cloud-sdk/bin/*

RUN gcloud components install kubectl \
    && rm -rf /google-cloud-sdk/.install/.backup

COPY estafette-extension-gke /
COPY templates /templates

RUN du -hd1 /google-cloud-sdk/.[^.]* /google-cloud-sdk/*
RUN du -hd1 /google-cloud-sdk/bin/.[^.]* /google-cloud-sdk/bin/*
RUN kubectl version

ENTRYPOINT ["/estafette-extension-gke"]