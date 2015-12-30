FROM golang:1.5.1
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN apt-get update && \
    apt-get -y install git sshpass

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

RUN git clone --depth=1 --recursive https://github.com/Nanocloud/nanocloud.git     && make -C nanocloud
RUN git clone --depth=1 --recursive https://github.com/Nanocloud/apps.git          && make -C apps
RUN git clone --depth=1 --recursive https://github.com/Nanocloud/iaas.git          && make -C iaas
RUN git clone --depth=1 --recursive https://github.com/Nanocloud/ldap.git          && make -C ldap
RUN git clone --depth=1 --recursive https://github.com/Nanocloud/history.git       && make -C history
RUN git clone --depth=1 --recursive https://github.com/Nanocloud/users.git         && make -C users

RUN mkdir -p nanocloud/bin/plugins/running

RUN mv iaas/bin/iaas nanocloud/bin/plugins/running/iaas && \
    mv users/bin/users nanocloud/bin/plugins/running/users && \
    mv ldap/bin/ldap nanocloud/bin/plugins/running/ldap && \
    mv history/bin/history nanocloud/bin/plugins/running/history && \
    mv apps/front/website/js/components/applications nanocloud/bin/front/js/components/applications && \
    mv apps/front/website/js/components/presenter nanocloud/bin/front/js/components/presenter && \
    mv iaas/front/website/js/components/services nanocloud/bin/front/js/components/services && \
    mv users/front/website/js/components/users nanocloud/bin/front/js/components/users && \
    mv history/front/website/js/components/history nanocloud/bin/front/js/components/history

RUN rm -rf iaas users ldap history apps && apt-get clean

ENV GOPATH="$GOPATH:/opt/build/backend"

EXPOSE 8081

WORKDIR /opt/build/backend/src/nanocloud.com/nanocloud/bin

CMD ["./nanocloud"]
