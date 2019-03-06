FROM golang:1.9.1

RUN echo "deb http://www.apache.org/dist/cassandra/debian 311x main" | tee -a /etc/apt/sources.list.d/cassandra.sources.list \
    && curl https://www.apache.org/dist/cassandra/KEYS | apt-key add - \
    && apt-get update && apt-get install cassandra -y

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

RUN cd /go/src/cassandra-web && go install

WORKDIR /go/src/cassandra-web

CMD ["cassandra-web"]