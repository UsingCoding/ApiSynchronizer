FROM debian:stable-slim

COPY bin/ /app/bin/

CMD /app/bin/apisynchronizer