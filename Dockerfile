FROM mosstech/nginx-with-prometheus:1.0
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
