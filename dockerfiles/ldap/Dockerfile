FROM golang:1.5.2
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com>

RUN go get -u github.com/constabulary/gb/...
RUN mkdir -p /build/

RUN git clone --depth=1 --recursive https://github.com/Nanocloud/ldap.git /build/ldap

RUN apt-get update && apt-get install libldap2-dev sshpass

WORKDIR /build/ldap
RUN gb build

RUN cp /build/ldap/src/ldap/script/*.sh /build/ldap/bin && chmod +x /build/ldap/bin/*.sh
CMD ["./bin/ldap"]
