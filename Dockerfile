FROM golang:1.23-alpine as builder

# Build jrpc-interceptor app
ADD . /source
WORKDIR /source
RUN go build -v .

# Copy jrpc-interceptor app to openresty (nginx) container
# openresty allows to run lua scripts
FROM openresty/openresty:alpine
ADD . /source
WORKDIR /source
COPY --from=builder /source/jrpc-interceptor .

# Syslog port
EXPOSE 514
# Prometheus metrics port
EXPOSE 9100

# Copy nginx template file
COPY nginx/nginx.conf /usr/local/openresty/nginx/conf/nginx.conf
RUN mkdir -p /var/log/nginx
RUN chmod 755 /var/log/nginx
COPY nginx/default.conf.template /etc/nginx/templates/nginx.conf.template

# add gettext for envsubst
RUN apk add --no-cache gettext

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]