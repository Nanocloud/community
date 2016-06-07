FROM node:5.6
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN npm install -g bower

RUN mkdir -p /opt/front
WORKDIR /opt/front

COPY package.json /tmp/package.json
RUN cd /tmp && npm install
RUN cp -a /tmp/node_modules /opt/front

COPY bower.json /tmp/bower.json
RUN cd /tmp && bower install --allow-root
RUN cp -a /tmp/bower_components /opt/front

COPY ./ /opt/front/

RUN ./node_modules/ember-cli/bin/ember build --environment=production

EXPOSE 8080
VOLUME ["/opt/front"]
