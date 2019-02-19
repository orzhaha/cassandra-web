FROM golang:1.9.1

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py \
    && python get-pip.py \
    && pip install cqlsh

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

RUN cd /go/src/cassandra-web && go install

WORKDIR /go/src/cassandra-web

CMD ["cassandra-web"]