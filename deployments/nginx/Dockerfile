FROM nginx:latest

ENV NGINX_PUBLIC /usr/share/nginx/html

COPY default.conf /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/nginx.conf

COPY images/. ${NGINX_PUBLIC}
