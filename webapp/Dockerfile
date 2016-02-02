FROM node:4.2.4
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN mkdir -p /opt/front
WORKDIR /opt/front
COPY ./ /opt/front/

RUN npm install && npm run setup

EXPOSE 8080
VOLUME ["/opt/front"]
