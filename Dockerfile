# build stage
FROM golang:1.9.1 AS build-env

RUN cd /go/src/ && git clone https://github.com/orzhaha/cassandra-web.git
RUN go get -u github.com/jteeuwen/go-bindata/...
RUN go-bindata -o=/go/src/cassandra-web/client/client.go -pkg=client /go/src/cassandra-web/client/dist/...
RUN cd /go/src/cassandra-web/service && go build 
 
# final stage
FROM golang:1.9.1

RUN echo "deb http://www.apache.org/dist/cassandra/debian 311x main" | tee -a /etc/apt/sources.list.d/cassandra.sources.list \
    && apt-get update && apt-get install curl gnupg -y \
    && curl https://www.apache.org/dist/cassandra/KEYS | apt-key add - \
    && apt-get update && apt-get install cassandra -y \
    && apt-get clean -y \ 
    && rm -fr /var/lib/apt/lists/* 

COPY --from=build-env /go/src/cassandra-web/service/service /
COPY --from=build-env /go/src/cassandra-web/service/config.yaml /

WORKDIR /

CMD ["./service"]
