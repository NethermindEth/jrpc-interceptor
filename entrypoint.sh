#!/bin/sh
./jrpc-interceptor -debug true &
envsubst '${SERVICE_TO_PROXY} ${LOG_SERVER_URL} ${LISTEN_PORT}' < /etc/nginx/templates/nginx.conf.template > /etc/nginx/conf.d/default.conf && exec nginx -g 'daemon off;'