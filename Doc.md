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

----

##### /rowtoken (get all row in table by token)

##### params

* table
* item
* prevnext
* pagesize

##### return：Json object

##### example：

```sh
curl \
  -X POST \
  -H 'Content-Type: application/json' \
  -d '{"table":"table_name", "item":"{field1:value1, field2:value2}", "prevnext": "next", "pagesize": 10}'
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

##### /save (eidt Save)

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
  -d '{"table":"table_name", "item":"{field1:value1, field2:value2}"}'
```

---

##### /describe (cqlsh describe)

##### params

* table

##### return：string

##### example：

```sh
curl \
  -X GET \
  http://localhost/describe?table=ca.abc
  -H 'Content-Type: application/json' \
```

---

##### /columns (get table columns)

##### params

* keyspace
* table

##### return：Json object

##### example：

```sh
curl \
  -X GET \
  http://localhost/columns \
  -H 'Content-Type: application/json' \
  -d '{"keyspace":"system_auth", "table":"roles"}'
```

---

##### /delete (delete table row)

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
  -d '{"table":"table_name", "item":"{field1:value1, field2:value2}"}'
```

---

##### /Find (find table row)

##### params

* table
* item
* order_by array[object{name:columns_name,order:asc|desc}] only clustering key and by order
* pagesize
* page

##### return：Json object

##### example：

```sh
curl \
  -X POST \
  http://localhost/save \
  -H 'Content-Type: application/json' \
  -d '{"table":"tablename","item":{"order_id":{"operator":"=","value":"123"},"order_time":{"operator":"=","value":"456"}}, "page": 1, "pagesize": 10}'
```

---


##### /Export (export file table data)

##### params

* table

##### return：file

##### example：

```sh
curl \
  -X POST \
  http://localhost/export?table=keyspace.table \
  -H 'Content-Type: application/json' \
```

---


##### /Import (import file table data)

##### form-data params

* table
* file

##### return：Josn object

##### example：

```sh
curl \
  -X POST \
  http://localhost/import \
  -F "table=keyspace.table" \
  -F "filecomment=This is an image file" \
  -F "file=@/home/user/Desktop/importfile" \
```

---