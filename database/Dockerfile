FROM postgres

WORKDIR /tmp/mta

# add a script to the docker entrypoint which will be executed when the
# container is run. This script will optionally setup a test database if
# the container is not running in prod.
COPY create-database.sh /docker-entrypoint-initdb.d/

RUN apt-get -qq update
RUN apt-get -qq install unzip

ADD http://web.mta.info/developers/data/nyct/subway/google_transit.zip /tmp/mta/

RUN unzip /tmp/mta/google_transit.zip
