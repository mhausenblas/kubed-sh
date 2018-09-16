FROM alpine:3.7
MAINTAINER Michael Hausenblas, mhausenblas.info

RUN apk add --no-cache curl

USER 1001

# install ttyd and kubed-sh
RUN curl -s -L -k https://github.com/tsl0922/ttyd/releases/download/1.4.2/ttyd_linux.x86_64 -o ttyd && \
    chmod 750 ttyd && \
    curl -s -L -k https://github.com/mhausenblas/kubed-sh/releases/download/0.5.1/kubed-sh-linux -o kubed-sh && \
    chmod 750 kubed-sh

EXPOSE 8080
CMD ["./ttyd -p 8080 ./kubed-sh"]
