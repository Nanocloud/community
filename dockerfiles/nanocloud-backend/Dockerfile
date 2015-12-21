FROM golang:1.5.1
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN apt-get update && \
    apt-get -y install git libldap2-dev sshpass

COPY repos_nanocloud /opt/build/backend/src/nanocloud.com

RUN curl -sL https://deb.nodesource.com/setup_5.x | bash - && apt-get install -y nodejs

RUN npm install -g \
    bower \
    grunt-cli \
    grunt \
    grunt-contrib-concat \
    grunt-contrib-copy \
    grunt-contrib-less \
    grunt-contrib-uglify \
    grunt-contrib-watch \
    grunt-ts \
    tslint \
    grunt-tslint \
    grunt-sync \
    tsd \
    http-server \
    concurrently \
    load-grunt-tasks

RUN echo '{"analytics": false}' > ~/.bowerrc

WORKDIR /opt/build/backend/src/nanocloud.com

RUN rm -rf nanocloud/.git && make -C nanocloud
RUN rm -rf apps/.git      && make -C apps
RUN rm -rf iaas/.git      && make -C iaas
RUN rm -rf ldap/.git      && make -C ldap
RUN rm -rf history/.git   && make -C history
RUN rm -rf users/.git     && make -C users

WORKDIR /opt/build/backend/src/nanocloud.com/nanocloud/

RUN mkdir -p bin/plugins/running

RUN ln -s ../../../../iaas/bin/iaas bin/plugins/running/iaas && \
    ln -s ../../../../users/bin/users bin/plugins/running/users && \
    ln -s ../../../../ldap/bin/ldap bin/plugins/running/ldap && \
    ln -s ../../../../history/bin/history bin/plugins/running/history && \
    ln -s ../../../../../apps/front/website/js/components/applications bin/front/js/components/applications && \
    ln -s ../../../../../apps/front/website/js/components/presenter bin/front/js/components/presenter && \
    ln -s ../../../../../iaas/front/website/js/components/services bin/front/js/components/services && \
    ln -s ../../../../../users/front/website/js/components/users bin/front/js/components/users && \
    ln -s ../../../../../history/front/website/js/components/history bin/front/js/components/history

ENV GOPATH="$GOPATH:/opt/build/backend"

EXPOSE 8081

WORKDIR /opt/build/backend/src/nanocloud.com/nanocloud/bin

CMD ["./nanocloud"]
