FROM al3xos/nginx-with-prometheus:1.17-3043616448-2026.11.09.0609
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
