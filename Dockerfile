FROM al3xos/nginx-with-prometheus:1.17-3072270697-2026.11.16.0609
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
