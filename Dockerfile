FROM alpine:3.18.3

ENV KUBECTL_VERSION="v1.25.13"

RUN apk add --update --upgrade --no-cache \
    curl \
  && rm -rf /var/cache/apk/*

RUN curl -L "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output /usr/bin/kubectl \
  && curl -LO "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl.sha256" \
  && echo "$(cat kubectl.sha256)  /usr/bin/kubectl" | sha256sum -c - \
  && chmod +x /usr/bin/kubectl \
  && kubectl version --client

FROM google/cloud-sdk:451.0.1-alpine

RUN apk add --update --upgrade --no-cache \
  && rm -rf google-cloud-sdk/bin/anthoscli \
  && rm -rf /var/cache/apk/*

RUN gcloud components install gke-gcloud-auth-plugin

LABEL maintainer="estafette.io" \
    description="The ${ESTAFETTE_GIT_NAME} component is an Estafette extension to deploy applications to a Kubernetes Engine cluster"

COPY --from=0 /usr/bin/kubectl /usr/bin/kubectl
COPY ${ESTAFETTE_GIT_NAME} /
COPY templates /templates

RUN mkdir -p ~/.kube

ENV ESTAFETTE_LOG_FORMAT="console" \
  GOOGLE_APPLICATION_CREDENTIALS="/key-file.json" \
  KUBECONFIG="/root/.kube/config"

ENTRYPOINT ["/${ESTAFETTE_GIT_NAME}"]
