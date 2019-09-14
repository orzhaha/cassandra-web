# cassandra-web

## Demo
![index](webui.gif)

## Features
* Table Row Prev Next Page
* Table Row Edit
* Table Row filter
* Table Row Delete
* Table Rown Find
* Table Definition
* Table Export
* Table Import
* CQL Query

---

## Usage

download
```
$ wget https://github.com/orzhaha/cassandra-web/releases/download/v1.0.3/linux.tar.gz
```

unzip
```
$ tar zxvf linux.tar.gz
```

npm install 
```
$ cd client && npm i && npm run build
```

run service
```
$ ./service -c config.yaml
```

#### depend

cqlsh 

---

## Docker

```sh
docker pull ipushc/cassandra-web
```
##### Environment

* HOST_PORT: ":80"
* CASSANDRA_HOST: cassandra host
* CASSANDRA_PORT: 9042
* CASSANDRA_USERNAME: username
* CASSANDRA_PASSWORD: password

---

## API

API [Doc](./Doc.md)

