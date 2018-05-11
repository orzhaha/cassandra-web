# cassandra-web
cassandra web ui

## API

#### /allkeyspace 取得所有KeySpace

##### 參數：無

##### 回傳：Json array

##### 範例：

```sh
http://localhost/allkeyspace
```

----

##### /alltablebykeyspace 取得KeySpace裡所有Table

##### 參數：

* keyspace

##### 回傳：Json array

##### 範例：

```sh
http://localhost/alltablebykeyspace?keyspace=abc
```

----

##### /allrowbytable 取得Table裡所有的Row

##### 參數：

* limit 筆數
* token_key 主鍵欄位名稱
* prev_token 上一頁最靠近的一筆資料
* next_token 下一頁最靠近的一筆資料

##### 回傳：Json object

##### 範例：

```sh
http://localhost/allrowbytable?table=ca.abc&limit=5&token_key=a&next_token=Serenity
```

---