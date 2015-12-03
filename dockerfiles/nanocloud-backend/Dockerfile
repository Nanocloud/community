FROM golang:1.4.3
MAINTAINER \
	Romain Soufflet <romain.soufflet@nanocloud.com> \
	Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN apt-get update && \
    apt-get -y install git libldap2-dev sshpass

COPY repos_nanocloud /opt/repos

RUN curl -sL https://deb.nodesource.com/setup_5.x | bash - && apt-get install -y nodejs

WORKDIR /opt/build
RUN git clone /opt/repos/front frontend
WORKDIR frontend
RUN npm install -g bower
RUN npm install
RUN bower install --allow-root
RUN npm run tsd install
RUN npm run grunt update-libs
RUN npm run grunt build

RUN mkdir -p /opt/build/backend/src/nanocloud.com
WORKDIR /opt/build/backend/src/nanocloud.com

RUN git clone /opt/repos/core core
RUN mkdir -p plugins

WORKDIR /opt/build/backend/src/nanocloud.com/plugins
RUN git clone /opt/repos/plugin_iaas iaas
RUN git clone /opt/repos/plugin_ldap ldap
RUN git clone /opt/repos/plugin_history history

WORKDIR /opt/build/backend/src/nanocloud.com/core
ENV GOPATH="$GOPATH:/opt/build/backend"
ENV NANOCONF="/opt/conf/config.json"
RUN make setup \
    && make haptic

EXPOSE 8081

CMD ["make", "serve"]
