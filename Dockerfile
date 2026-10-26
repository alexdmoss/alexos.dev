FROM al3xos/nginx-with-prometheus:1.17@sha256:f4c9da28c3687e801c62083caf749a4f87b7861a0a72f64100a0138087e5cd6b
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
