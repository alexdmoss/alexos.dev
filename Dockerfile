FROM nginx-with-prometheus:latest
COPY config/default.conf /etc/nginx/conf.d/
COPY www/ /app/
