FROM al3xos/nginx-with-prometheus:1.17-3053395545-2026.11.11.0827
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
