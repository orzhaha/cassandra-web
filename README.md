# cassandra-web

[![latest-tag](https://badgen.net/github/release/orzhaha/cassandra-web)](https://github.com/orzhaha/cassandra-web/releases)


## Demo
![index](webui.gif)

## Features
* Theme Dark
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

## Supported Cassandra Versions
------------------

 2.1.x | 2.2.x | 3.x.x
 -------| ------| ---------
 yes | yes | yes



---

## Usage

download
```
$ wget https://github.com/orzhaha/cassandra-web/releases/download/v1.0.5/linux.tar.gz
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
* READ_ONLY: false (default)
* CASSANDRA_HOST: cassandra host
* CASSANDRA_PORT: 9042
* CASSANDRA_USERNAME: username
* CASSANDRA_PASSWORD: password
* APP_PATH: path

---

## API

API [Doc](./Doc.md)

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=orzhaha/cassandra-web&type=Date)](https://star-history.com/#orzhaha/cassandra-web&Date)

