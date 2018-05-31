FROM golang:1.9.1

RUN cd /go/src

RUN git clone https://github.com/orzhaha/cassandra-web.git

RUN cd cassandra-web && go install

WORKDIR /go/src/cassandra-web

CMD ["cassandra-web"]