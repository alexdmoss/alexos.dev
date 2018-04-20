FROM nginx:alpine
COPY config/nginx/site.conf /etc/nginx/conf.d/default.conf
COPY app/ /app/
EXPOSE 32080
WORKDIR "/app"
