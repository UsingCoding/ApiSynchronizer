FROM debian:stable-slim

RUN apt-get update && apt-get install -y git

RUN apt clean autoclean && \
    apt autoremove -y && \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

COPY bin/ /app/bin/

ENTRYPOINT [ "/app/bin/apisynchronizer" ]