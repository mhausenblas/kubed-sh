FROM alpine:3.7
MAINTAINER Michael Hausenblas, mhausenblas.info

RUN apk add --no-cache curl

USER 1001

ENV KUBECTL_BINARY=/tmp/ko/kubectl

# install ttyd and kubed-sh
RUN mkdir /tmp/ko && curl -s -L -k https://github.com/tsl0922/ttyd/releases/download/1.4.2/ttyd_linux.x86_64 -o /tmp/ko/ttyd && \
    chmod 750 /tmp/ko/ttyd && \
    curl -s -L -k https://github.com/mhausenblas/kubed-sh/releases/download/0.5.1/kubed-sh-linux -o /tmp/ko/kubed-sh && \
    chmod 750 /tmp/ko/kubed-sh && \
    curl -s -L -k https://storage.googleapis.com/kubernetes-release/release/$(curl -s -k https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl -o /tmp/ko/kubectl && \
    chmod 750 /tmp/ko/kubectl

EXPOSE 8080
ENTRYPOINT [ "/tmp/ko/ttyd", "-p 8080", "/tmp/ko/kubed-sh" ]
