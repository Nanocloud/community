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

RUN git clone --recursive /opt/repos/iaas iaas
RUN make -C iaas

RUN git clone --recursive /opt/repos/ldap ldap
RUN make -C ldap

RUN git clone --recursive /opt/repos/history history
RUN make -C history

RUN git clone --recursive /opt/repos/users users
RUN make -C users

WORKDIR /opt/build/backend/src/nanocloud.com/nanocloud/bin

RUN mkdir -p plugins/running

RUN ln -s ../../../../iaas/bin/iaas plugins/running/iaas
RUN ln -s ../../../../ldap/bin/ldap plugins/running/ldap
RUN ln -s ../../../../history/bin/history plugins/running/history
RUN ln -s ../../../../users/bin/users plugins/running/users

RUN ln -s ../../../../../iaas/front/website/js/components/iaas front/js/components/iaas
RUN ln -s ../../../../../history/front/website/js/components/history front/js/components/history
RUN ln -s ../../../../../users/front/website/js/components/users front/js/components/users

ENV GOPATH="$GOPATH:/opt/build/backend"
ENV NANOCONF="/opt/conf/config.json"

EXPOSE 8081

CMD ["./nanocloud"]
