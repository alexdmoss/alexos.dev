FROM al3xos/nginx-with-prometheus:1.15
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
