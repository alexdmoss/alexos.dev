FROM al3xos/nginx-with-prometheus:1.17@sha256:085178b5fc011833ef5ff0765d5d982824197b78eecd050b841f8adde8510aaf
USER nginx
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
