FROM golang:1.4.3
MAINTAINER \
	Romain Soufflet <romain.soufflet@nanocloud.com> \
	Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN apt-get update && \
    apt-get -y install git libldap2-dev sshpass

RUN curl -sL https://deb.nodesource.com/setup_5.x | bash - && apt-get install -y nodejs

RUN mkdir /opt/build
WORKDIR /opt/build
RUN git clone https://github.com/Gentux/users.git && \
	cd users && \
	git submodule update --init && \
	make initial && \
	make

RUN git clone https://github.com/Nanocloud/iaas.git && \
	cd iaas && \
	git submodule update --init && \
	make initial && \
	make

RUN git clone https://github.com/Nanocloud/ldap.git && \
	cd ldap && \
	git submodule update --init && \
	make initial && \
	make

RUN git clone https://github.com/Nanocloud/history.git && \
	cd history && \
	git submodule update --init && \
	make initial && \
	make

RUN git clone https://github.com/Nanocloud/oauth.git && \
	cd oauth && \
	git submodule update --init && \
	make initial && \
	make

RUN git clone https://github.com/Nanocloud/apps.git && \
	cd apps && \
	git submodule update --init && \
	make initial && \
	make

RUN git clone https://github.com/Nanocloud/nanocloud.git && \
	cd nanocloud && \
	git submodule update --init && \
	make initial && \
	make

EXPOSE 8081
CMD ["./nanocloud/nanocloud"]
