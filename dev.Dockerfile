FROM golang:1.9.1

RUN echo "deb http://www.apache.org/dist/cassandra/debian 311x main" | tee -a /etc/apt/sources.list.d/cassandra.sources.list \
    && curl https://www.apache.org/dist/cassandra/KEYS | apt-key add - \
    && apt-get update && apt-get install cassandra -y

RUN wget https://nodejs.org/dist/v8.5.0/node-v8.5.0-linux-x64.tar.gz \
    && tar -xf node-v8.5.0-linux-x64.tar.gz --directory /usr/local --strip-components 1 \
    && rm -rf node-v8.5.0-linux-x64.tar.gz

CMD ["tail", "-f", "/dev/null"]