FROM google/cloud-sdk:222.0.0-alpine

LABEL maintainer="estafette.io" \
      description="The estafette-extension-gke component is an Estafette extension to deploy applications to a Kubernetes Engine cluster"

ARG KUBECTL_VERSION=1.10.7
ENV KUBECTL_VERSION=$KUBECTL_VERSION

RUN curl https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl --output /google-cloud-sdk/bin/kubectl \
    && chmod +x /google-cloud-sdk/bin/kubectl

COPY estafette-extension-gke /
COPY templates /templates

RUN kubectl version --client
RUN type -a kubectl

ENTRYPOINT ["/estafette-extension-gke"]