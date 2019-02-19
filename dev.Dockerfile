FROM golang:1.9.1

RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py \
    && python get-pip.py \
    && pip install cqlsh

RUN wget https://nodejs.org/dist/v8.5.0/node-v8.5.0-linux-x64.tar.gz \
    && tar -xf node-v8.5.0-linux-x64.tar.gz --directory /usr/local --strip-components 1 \
    && rm -rf node-v8.5.0-linux-x64.tar.gz

CMD ["tail", "-f", "/dev/null"]