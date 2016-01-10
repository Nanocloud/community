FROM alpine:3.3
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN apk add --no-cache \
	git \
	openssh-client \
	perl

RUN git clone --depth=1 --recursive https://github.com/Nanocloud/community.git

VOLUME ["/var/lib/nanocloud"]

COPY entrypoint.sh /entrypoint.sh

CMD ["./entrypoint.sh"]
