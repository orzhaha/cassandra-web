# build stage
FROM ipushc/golangxnode:1.9.1-v12 AS build-env

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git

RUN cd /go/src/cassandra-web/client && npm i && npm run build

RUN go get -u github.com/jteeuwen/go-bindata/... \
    && cd /go/src/cassandra-web/ \
    && go-bindata -o=client/client.go -pkg=client client/dist/...

RUN cd /go/src/cassandra-web/service && go build 

# final stage
FROM debian:stable

RUN apt-get update && apt-get install --no-install-recommends curl gnupg -y \
    && echo "deb http://www.apache.org/dist/cassandra/debian 311x main" | tee -a /etc/apt/sources.list.d/cassandra.sources.list \
    && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys A278B781FE4B2BDA \
    # && curl https://www.apache.org/dist/cassandra/KEYS | apt-key add - \
    && apt-get update && apt-get install cassandra -y \
    && apt-get clean -y \
    && rm -fr /var/lib/apt/lists/* 

COPY --from=build-env /go/src/cassandra-web/service/service /
COPY --from=build-env /go/src/cassandra-web/service/config.yaml /

WORKDIR /

CMD ["./service"]
