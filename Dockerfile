FROM google/cloud-sdk:234.0.0-alpine

LABEL maintainer="estafette.io" \
      description="The estafette-extension-gke component is an Estafette extension to deploy applications to a Kubernetes Engine cluster"

RUN curl https://storage.googleapis.com/kubernetes-release/release/v1.10.7/bin/linux/amd64/kubectl --output /google-cloud-sdk/bin/kubectl \
    && chmod +x /google-cloud-sdk/bin/kubectl \
    && kubectl version --client

COPY estafette-extension-gke /
COPY templates /templates

ENTRYPOINT ["/estafette-extension-gke"]