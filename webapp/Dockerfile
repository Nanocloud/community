FROM node:4.2.4
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN mkdir -p /opt/front
WORKDIR /opt/front
COPY ./ /opt/front/

RUN npm install && npm run setup
RUN cd modules/apps && npm install && npm run setup && mv website/js/components/applications ../../website/js/components && mv website/js/components/presenter ../../website/js/components
RUN cd modules/iaas && npm install && npm run setup && mv website/js/components/services ../../website/js/components
RUN cd modules/users && npm install && npm run setup && mv website/js/components/users ../../website/js/components
RUN cd modules/history && npm install && npm run setup && mv website/js/components/history ../../website/js/components && mv website/ts/components/history/views ../../website/js/components/history

EXPOSE 8080
VOLUME ["/opt/front"]