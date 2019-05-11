FROM nginx:alpine
COPY config/nginx.conf /etc/nginx/conf.d/default.conf
COPY www/ /app/
EXPOSE 80
WORKDIR /app
