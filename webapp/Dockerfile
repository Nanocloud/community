FROM node:5.6
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN npm install -g bower ember-cli

RUN mkdir -p /opt/front
WORKDIR /opt/front
COPY ./ /opt/front/

ADD package.json /tmp/package.json
RUN cd /tmp && npm install

RUN cp -r /tmp/node_modules /opt/front/
RUN bower install --allow-root
RUN ember build --environment="production"

EXPOSE 8080
VOLUME ["/opt/front"]
