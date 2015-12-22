FROM golang:1.5.1
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN apt-get update && \
    apt-get -y install git libldap2-dev sshpass

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

RUN mkdir -p /opt/build/backend/src/nanocloud.com
WORKDIR /opt/build/backend/src/nanocloud.com

RUN git clone --depth=1 https://github.com/Nanocloud/nanocloud.git    && rm -rf nanocloud/.git && make -C nanocloud
RUN git clone --depth=1 https://github.com/Nanocloud/apps.git         && rm -rf apps/.git      && make -C apps
RUN git clone --depth=1 https://github.com/Nanocloud/iaas.git         && rm -rf iaas/.git      && make -C iaas
RUN git clone --depth=1 https://github.com/Nanocloud/ldap.git         && rm -rf ldap/.git      && make -C ldap
RUN git clone --depth=1 https://github.com/Nanocloud/history.git      && rm -rf history/.git   && make -C history
RUN git clone --depth=1 https://github.com/Nanocloud/users.git        && rm -rf users/.git     && make -C users

RUN mkdir -p nanocloud/bin/plugins/running

RUN cp iaas/bin/iaas bin/plugins/running/iaas && \
    cp users/bin/users bin/plugins/running/users && \
    cp ldap/bin/ldap bin/plugins/running/ldap && \
    cp history/bin/history bin/plugins/running/history && \
    cp apps/front/website/js/components/applications bin/front/js/components/applications && \
    cp apps/front/website/js/components/presenter bin/front/js/components/presenter && \
    cp iaas/front/website/js/components/services bin/front/js/components/services && \
    cp users/front/website/js/components/users bin/front/js/components/users && \
    cp history/front/website/js/components/history bin/front/js/components/history

RUN rm -rf iaas users ldap history nanocloud apps && apt-get clean

ENV GOPATH="$GOPATH:/opt/build/backend"

EXPOSE 8081

WORKDIR /opt/build/backend/src/nanocloud.com/nanocloud/bin

CMD ["./nanocloud"]
