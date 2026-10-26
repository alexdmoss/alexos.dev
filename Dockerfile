FROM al3xos/nginx-with-prometheus:1.17@sha256:940fd9661e09532fa5ba9aa124d16571232fee315f5411c728dc4143afb8d2ab
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
