# cassandra-web
cassandra web ui

## Docker

```sh
docker pull ipushc/cassandra-web
```

----

## 環境變數

* HOST_PORT: ":80"
* CASSANDRA_HOST: cassandra
* CASSANDRA_PORT: 9042

----

## API

#### /keyspace 取得所有KeySpace

##### 參數：無

##### 回傳：Json array

##### 範例：

```sh
curl \
  -X GET \
  -H 'Content-Type: application/json' \
  http://localhost/keyspace
```

----

##### /table 取得KeySpace裡所有Table

##### 參數：

* keyspace

##### 回傳：Json array

##### 範例：

```sh
curl \
  -X GET \
  -H 'Content-Type: application/json' \
  http://localhost/table?keyspace=abc
```

----

##### /row 取得Table裡所有的Row

##### 參數：

* limit 筆數
* token_key 主鍵欄位名稱
* prev_token 上一頁最靠近的一筆資料
* next_token 下一頁最靠近的一筆資料

##### 回傳：Json object

##### 範例：

```sh
curl \
  -X GET \
  -H 'Content-Type: application/json' \
  http://localhost/row?table=ca.abc&limit=5&token_key=a&next_token=Serenity
```

---

##### /query 執行cql query

##### 參數：

* query cql query

##### 回傳：Json object

##### 範例：

```sh
curl \
  -X POST \
  http://localhost/query \
  -H 'Content-Type: application/json' \
  -d '{"query":"cql query"}'
```

---