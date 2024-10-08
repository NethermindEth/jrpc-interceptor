#!/bin/sh
./jrpc-interceptor -debug=${LOG_SERVER_DEBUG:-true} -listenSyslog=${LOG_SERVER_URL:-"0.0.0.0:514"} -listenHTTP=${PROMETHEUS_URL:-"0.0.0.0:9100"} -usePrometheus=${USE_PROMETHEUS:-true}  &
envsubst '${SERVICE_TO_PROXY} ${LOG_SERVER_URL} ${LISTEN_PORT}' < /etc/nginx/templates/nginx.conf.template > /etc/nginx/conf.d/default.conf && exec nginx -g 'daemon off;'