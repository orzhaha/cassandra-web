FROM golang:1.9.1

RUN cd /go/src && \
    git clone https://github.com/orzhaha/cassandra-web.git && \
    cd cassandra-web && \
    go build -o cassandra-web

CMD ["./cassandra-web"]