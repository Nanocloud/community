FROM node:4.2.4
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN mkdir -p /opt/canva
WORKDIR /opt/canva
COPY ./ /opt/canva/

RUN npm install -g bower ember-cli && npm install
RUN bower install --allow-root
RUN ember build --environment="production"

EXPOSE 8080
VOLUME ["/opt/canva"]
