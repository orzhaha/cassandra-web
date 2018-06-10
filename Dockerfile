FROM ipushc/observedcat

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

RUN cd /go/src/cassandra-web && go install

WORKDIR /go/src/cassandra-web

CMD ["cassandra-web"]