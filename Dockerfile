FROM al3xos/nginx-with-prometheus:1.17-3016388354-2026.11.02.0609
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
