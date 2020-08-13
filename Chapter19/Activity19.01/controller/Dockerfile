# FROM registry.fedoraproject.org/fedora:32 as builder
FROM fedora:32 as builder

LABEL version="1.0" \
      maintainer="fmasood@redhat.com"

ENV GOPATH /root/go

RUN yum -y update && yum -y install golang git dep
RUN mkdir -p /root/go/src/k8s.io && cd /root/go/src/k8s.io && \
    git clone --single-branch --branch release-1.14 https://github.com/kubernetes/code-generator.git &&  \
    git clone --single-branch --branch release-1.14  https://github.com/kubernetes/apimachinery.git && \
    mkdir -p /root/go/src/github.com/example-inc/pod-normaliser-controller


WORKDIR /root/go/src/github.com/example-inc/pod-normaliser-controller
COPY ./*.go ./Gopkg.toml /root/go/src/github.com/example-inc/pod-normaliser-controller/
ADD pkg /root/go/src/github.com/example-inc/pod-normaliser-controller/pkg


ENV ROOT_PACKAGE github.com/example-inc/pod-normaliser-controller
ENV CUSTOM_RESOURCE_NAME podlifecycleconfig
ENV CUSTOM_RESOURCE_VERSION v1

RUN dep ensure
RUN cd /root/go/src/k8s.io/code-generator && ./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
#
#
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
#


#FROM registry.fedoraproject.org/fedora-minimal:32
FROM fedora:32
RUN groupadd appgroup && useradd appuser -G appgroup
COPY --from=builder /root/go/src/github.com/example-inc/pod-normaliser-controller/main /app/
WORKDIR /app
USER appuser
CMD ["./main"]
