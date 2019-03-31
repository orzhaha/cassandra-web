# cassandra-web
cassandra web ui
![index](webui.gif)

## Features
* Table Row Prev Next Page
* Table Row Edit
* Table Row Delete
* Table Rown Find
* Table Definition
* Table Export
* Table Import
* CQL Query

## Docker

```sh
docker pull ipushc/cassandra-web
```

----

## Environment

* HOST_PORT: ":80"
* CASSANDRA_HOST: cassandra host
* CASSANDRA_PORT: 9042
* CASSANDRA_USERNAME: username
* CASSANDRA_PASSWORD: password

----

## API

API [Doc](./Doc.md)

## Issue

* CQL Data Types Map. JSON only allows key names to be strings.
* JSON int64 to string