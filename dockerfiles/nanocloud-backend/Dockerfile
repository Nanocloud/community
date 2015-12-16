FROM golang:1.4.3
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN apt-get update && \
    apt-get -y install git libldap2-dev sshpass

COPY repos_nanocloud /opt/repos

RUN curl -sL https://deb.nodesource.com/setup_5.x | bash - && apt-get install -y nodejs

RUN mkdir -p /opt/build/backend/src/nanocloud.com
WORKDIR /opt/build/backend/src/nanocloud.com

RUN echo '{"analytics": false}' > ~/.bowerrc

RUN git clone --recursive /opt/repos/nanocloud nanocloud
RUN make -C nanocloud

RUN git clone --recursive /opt/repos/apps apps
RUN make -C apps

RUN git clone --recursive /opt/repos/iaas iaas
RUN make -C iaas

RUN git clone --recursive /opt/repos/ldap ldap
RUN make -C ldap

RUN git clone --recursive /opt/repos/history history
RUN make -C history

RUN git clone --recursive /opt/repos/users users
RUN make -C users

WORKDIR /opt/build/backend/src/nanocloud.com/nanocloud/

RUN mkdir -p bin/plugins/running

RUN ln -s ../../../../apps/bin/apps bin/plugins/running/apps
RUN ln -s ../../../../iaas/bin/iaas bin/plugins/running/iaas
RUN ln -s ../../../../users/bin/users bin/plugins/running/users
RUN ln -s ../../../../ldap/bin/ldap bin/plugins/running/ldap
RUN ln -s ../../../../history/bin/history bin/plugins/running/history

RUN ln -s ../../../../../apps/front/website/js/components/applications bin/front/js/components/applications
RUN ln -s ../../../../../apps/front/website/js/components/presenter bin/front/js/components/presenter
RUN ln -s ../../../../../iaas/front/website/js/components/services bin/front/js/components/services
RUN ln -s ../../../../../users/front/website/js/components/users bin/front/js/components/users
RUN ln -s ../../../../../history/front/website/js/components/history bin/front/js/components/history

ENV GOPATH="$GOPATH:/opt/build/backend"

EXPOSE 8081

WORKDIR /opt/build/backend/src/nanocloud.com/nanocloud/bin

CMD ["./nanocloud"]
