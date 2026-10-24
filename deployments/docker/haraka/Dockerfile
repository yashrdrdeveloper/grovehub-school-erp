FROM node:22-alpine

ARG HARAKA_VERSION=3.3.5
ARG HARAKA_RELAY_VERSION=1.0.2
ARG HARAKA_DKIM_VERSION=1.3.1

RUN apk add --no-cache openssl \
    && npm install --global "Haraka@${HARAKA_VERSION}" \
    && haraka -i /opt/haraka \
    && npm install --prefix /opt/haraka \
        "haraka-plugin-relay@${HARAKA_RELAY_VERSION}" \
        "haraka-plugin-dkim@${HARAKA_DKIM_VERSION}" \
    && mkdir -p /opt/haraka/queue /opt/haraka/config/dkim

COPY deployments/docker/haraka/config/ /opt/haraka/config/
COPY deployments/docker/haraka/entrypoint.sh /usr/local/bin/grovehub-haraka-entrypoint
RUN chmod +x /usr/local/bin/grovehub-haraka-entrypoint

EXPOSE 2525

ENTRYPOINT ["/usr/local/bin/grovehub-haraka-entrypoint"]
