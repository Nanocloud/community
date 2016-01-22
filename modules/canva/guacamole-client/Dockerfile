FROM tomcat:8.0.20-jre7
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN apt-get update && apt-get -y install maven \
    openjdk-7-jdk

ENV GUAC_VERSION=0.9.8
ENV GUACAMOLE_HOME="/etc/guacamole"

RUN wget http://sourceforge.net/projects/guacamole/files/current/binary/guacamole-${GUAC_VERSION}.war/download -O /usr/local/tomcat/webapps/guacamole.war

RUN mkdir -p /opt/noauth-logged
COPY ./noauth-logged /opt/noauth-logged
WORKDIR /opt/noauth-logged

RUN mvn package
RUN mkdir -p $GUACAMOLE_HOME/extensions/ && cp ./target/guacamole-auth-noauthlogged-${GUAC_VERSION}.jar $GUACAMOLE_HOME/extensions/

RUN mvn clean

ENV JAVA_OPTS="-Djava.library.path=/usr/local/apr/lib -Djava.security.egd=file:/dev/./urandom -Djava.awt.headless=true -Xmx1024m -XX:MaxPermSize=512m -XX:+UseConcMarkSweepGC"

EXPOSE 8080

CMD ["catalina.sh", "run"]
