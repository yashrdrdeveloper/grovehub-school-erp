FROM node:22-alpine

ARG HARAKA_VERSION=3.3.4
ARG HARAKA_RELAY_VERSION=1.0.2

RUN npm install --global "Haraka@${HARAKA_VERSION}" \
    && haraka -i /opt/haraka \
    && npm install --prefix /opt/haraka "haraka-plugin-relay@${HARAKA_RELAY_VERSION}" \
    && mkdir -p /opt/haraka/queue

COPY deployments/docker/haraka/config/ /opt/haraka/config/

EXPOSE 2525

CMD ["haraka", "-c", "/opt/haraka"]
