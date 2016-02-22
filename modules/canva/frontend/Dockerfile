FROM node:4.2.4
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN mkdir -p /opt/canva
WORKDIR /opt/canva
COPY ./ /opt/canva/

RUN npm install && npm install -g bower
RUN bower install --allow-root
RUN npm run build

EXPOSE 8080
VOLUME ["/opt/canva"]
