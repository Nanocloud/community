FROM nginx:1.9
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

COPY ./conf/nginx.conf /etc/nginx/conf.d/default.conf
COPY ./conf/certificates /etc/nginx/ssl/

EXPOSE 80
EXPOSE 443
