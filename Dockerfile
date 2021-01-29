# build client stage
FROM ipushc/golangxnode:1.9.1-v12 AS build-client-env

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

RUN cd /go/src/cassandra-web/client && npm i && npm run build

# build server stage
FROM golang:1.9.1-alpine AS build-server-env

RUN apk add --no-cache git

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

ENV GO111MODULE=on

COPY --from=build-client-env /go/src/cassandra-web/client/dist /go/src/cassandra-web/client/dist

RUN go get -u github.com/jteeuwen/go-bindata/... \
    && cd /go/src/cassandra-web/ \
    && go-bindata -o=client/client.go -pkg=client client/dist/...

RUN cd /go/src/cassandra-web/service && go build


# final stage
FROM alpine:3.13.1

RUN wget https://downloads.datastax.com/enterprise/cqlsh-astra.tar.gz \
    && tar zxvf cqlsh-astra.tar.gz \
    && mv cqlsh-astra/bin/cqlsh sbin/cqlsh \
    && mv cqlsh-astra/bin/cqlsh.py sbin/cqlsh.py \
    && mv cqlsh-astra/bin/dsecqlsh.py sbin/dsecqlsh.py \
    && mv cqlsh-astra/pylib/ / \
    && mv cqlsh-astra/zipfiles/ / \
    && apk add --no-cache python2

COPY --from=build-server-env /go/src/cassandra-web/service/service /
COPY --from=build-client-env /go/src/cassandra-web/service/config.yaml /

WORKDIR /

CMD ["./service"]
