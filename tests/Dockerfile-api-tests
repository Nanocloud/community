FROM node:0.10.42
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN npm install -g newman@1
COPY ./ /opt

WORKDIR /opt

CMD newman -x --insecure -c /opt/api/APItests.postman_collection -e /opt/api/NanoEnv.postman_environment -r 120000
