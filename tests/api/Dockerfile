FROM node:5.6
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

WORKDIR /opt
RUN npm install -g mocha
COPY ./ /opt
RUN rm -rf ./node_modules && npm install

CMD ["mocha", "index.js"]
