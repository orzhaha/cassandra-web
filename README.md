# cassandra-web
cassandra web ui

## Docker

```sh
docker pull ipushc/cassandra-web
```

----

## Environment

* HOST_PORT: ":80"
* CASSANDRA_HOST: cassandra
* CASSANDRA_PORT: 9042

----

## API

#### /keyspace (get all keyspace)

##### params

##### return：Json array

##### example：

```sh
curl \
  -X GET \
  -H 'Content-Type: application/json' \
  http://localhost/keyspace
```

----

##### /table (get all table in keyspace)

##### params

* keyspace

##### return：Json array

##### example：

```sh
curl \
  -X GET \
  -H 'Content-Type: application/json' \
  http://localhost/table?keyspace=abc
```

----

##### /row (get all row in table)

##### params

* table
* page
* pagesize

##### return：Json object

##### example：

```sh
curl \
  -X GET \
  -H 'Content-Type: application/json' \
  http://localhost/row?table=ca.abc&limit=5&token_key=a&next_token=Serenity
```

---

##### /query (cql query)

##### params

* query cql query

##### return：Json object

##### example：

```sh
curl \
  -X POST \
  http://localhost/query \
  -H 'Content-Type: application/json' \
  -d '{"query":"cql query"}'
```

---

##### /Save eidt Save

##### params

* table
* item

##### return：Json object

##### example：

```sh
curl \
  -X POST \
  http://localhost/save \
  -H 'Content-Type: application/json' \
  -d '{"table":"table_name", "item":"{field1:value1, field2:value2}"}'
```

---

## Issue

* CQL Data Types Map. JSON only allows key names to be strings.
* JSON int64 to string