# build client stage
FROM ipushc/golangxnode:1.16-v12 AS build-client-env

WORKDIR /workspace

COPY client/ client/
RUN cd /workspace/client && npm i && npm run build

# build server stage
FROM golang:1.16-alpine AS build-server-env

WORKDIR /workspace

RUN apk add --no-cache git

COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/
COPY service/ service/

RUN cd service && GO111MODULE=on go build -mod vendor

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

COPY --from=build-server-env /workspace/service/service /
COPY --from=build-server-env /workspace/service/config.yaml /
COPY --from=build-client-env /workspace/client/dist /client/dist

WORKDIR /

CMD ["./service"]
