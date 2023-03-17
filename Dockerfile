# build client stage
FROM ipushc/golangxnode:1.16-v12 AS build-client-env

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

RUN cd /go/src/cassandra-web/client && npm i && npm run build

# build server stage
FROM golang:1.20.2-alpine AS build-server-env

RUN apk add --no-cache git

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

ENV GO111MODULE=on

RUN cd /go/src/cassandra-web/service && go build -mod vendor


# final stage
FROM alpine:3.13.10

RUN wget https://downloads.datastax.com/enterprise/cqlsh-astra.tar.gz \
    && tar zxvf cqlsh-astra.tar.gz \
    && mv cqlsh-astra/bin/cqlsh sbin/cqlsh \
    && mv cqlsh-astra/bin/cqlsh.py sbin/cqlsh.py \
    && mv cqlsh-astra/bin/dsecqlsh.py sbin/dsecqlsh.py \
    && mv cqlsh-astra/pylib/ / \
    && mv cqlsh-astra/zipfiles/ / \
    && apk add --no-cache python2

RUN wget https://github.com/masumsoft/cassandra-exporter/releases/download/v1.0.4/cassandra-exporter-linux.zip \
    && apk add zip \
    && unzip cassandra-exporter-linux.zip \
    && mv cassandra-exporter-linux/export-linux /sbin/cexport \
    && mv cassandra-exporter-linux/import-linux /sbin/cimport

RUN apk add gcompat \
    && apk add build-base \
    && apk upgrade zlib expat

COPY --from=build-server-env /go/src/cassandra-web/service/service /
COPY --from=build-client-env /go/src/cassandra-web/service/config.yaml /
COPY --from=build-client-env /go/src/cassandra-web/client/dist /client/dist

RUN adduser -D nonroot 
ENV HOME /home/nonroot
USER nonroot

WORKDIR /

CMD ["./service"]
