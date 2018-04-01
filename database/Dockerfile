FROM postgres

RUN apt-get -qq update
RUN apt-get -qq install unzip

ADD http://web.mta.info/developers/data/nyct/subway/google_transit.zip /tmp/mta/

RUN unzip -o /tmp/mta/google_transit.zip
