FROM tomcat:8.0.20-jre7
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN apt-get update && apt-get -y install maven \
    openjdk-7-jdk \
    git

ENV GUAC_VERSION=0.9.8
ENV GUACAMOLE_HOME="/etc/guacamole"

RUN wget http://sourceforge.net/projects/guacamole/files/current/binary/guacamole-${GUAC_VERSION}.war/download -O /usr/local/tomcat/webapps/guacamole.war

WORKDIR /opt/
RUN git clone https://github.com/Nanocloud/noauth-logged.git

RUN mvn package -f noauth-logged/pom.xml
RUN mkdir -p $GUACAMOLE_HOME/extensions/ && cp ./noauth-logged/target/guacamole-auth-noauthlogged-${GUAC_VERSION}.jar $GUACAMOLE_HOME/extensions/

RUN mvn clean -f noauth-logged/pom.xml

ENV JAVA_OPTS="-Djava.library.path=/usr/local/apr/lib -Djava.security.egd=file:/dev/./urandom -Djava.awt.headless=true -Xmx1024m -XX:MaxPermSize=512m -XX:+UseConcMarkSweepGC"

EXPOSE 8080

CMD ["catalina.sh", "run"]
