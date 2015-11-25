FROM tomcat:8.0.20-jre7
MAINTAINER Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN apt-get update && apt-get -y install maven \
    openjdk-7-jdk \
    git

ENV GUAC_VERSION=0.9.8
ENV GUACAMOLE_HOME="/etc/guacamole"

RUN mkdir -p $GUACAMOLE_HOME/extensions

COPY repos_guacamole /opt/repos

WORKDIR /opt/build/guacamole-client-${GUAC_VERSION}

RUN curl -fsSL http://sourceforge.net/projects/guacamole/files/current/source/guacamole-client-$GUAC_VERSION.tar.gz/download | tar xzf - -C /opt/build \
    && mvn package \
    && cp ./guacamole/target/guacamole-${GUAC_VERSION}.war /usr/local/tomcat/webapps/guacamole.war
RUN git clone /opt/repos/noauth-logged ./extensions/noauth-logged
RUN mvn package -f ./extensions/noauth-logged/pom.xml
RUN cp ./extensions/noauth-logged/target/guacamole-auth-noauthlogged-${GUAC_VERSION}.jar $GUACAMOLE_HOME/extensions/

ENV JAVA_OPTS="-Djava.library.path=/usr/local/apr/lib -Djava.security.egd=file:/dev/./urandom -Djava.awt.headless=true -Xmx1024m -XX:MaxPermSize=512m -XX:+UseConcMarkSweepGC"

EXPOSE 8080

CMD ["catalina.sh", "run"]
