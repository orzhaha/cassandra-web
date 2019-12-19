package main

import (
	"cassandra-web/client"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gocql/gocql"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

const (
	PartitionKey  = "partition_key"
	ClusteringKey = "clustering"
)

// init 初始化
func init() {
	// 反序列化float64精准度問題處理
	decodeNumberAsInt64IfPossible := func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {
		case jsoniter.NumberValue:
			var number json.Number
			iter.ReadVal(&number)
			i, err := strconv.ParseInt(string(number), 10, 64)
			if err == nil {
				*(*interface{})(ptr) = i
				return
			}
			f, err := strconv.ParseFloat(string(number), 64)
			if err == nil {
				*(*interface{})(ptr) = f
				return
			}
			// Not much we can do here.
		default:
			*(*interface{})(ptr) = iter.Read()
		}
	}
	jsoniter.RegisterTypeDecoderFunc("interface {}", decodeNumberAsInt64IfPossible)

	// 設置需要釋放的目錄
	dirs := []string{"client"}
	for _, dir := range dirs {
		// 解壓dir目錄到當前目錄
		if err := client.RestoreAssets("./", dir); err != nil {
			break
		}
	}
}

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	SystemSchemaKey = "system_schema"
)

var env envStruct

// envStruct type
type envStruct struct {
	HostPort          string `mapstructure:"HOST_PORT" json:"HOST_PORT"`
	CassandraHost     string `mapstructure:"CASSANDRA_HOST" json:"CASSANDRA_HOST"`
	CassandraPort     int    `mapstructure:"CASSANDRA_PORT" json:"CASSANDRA_PORT"`
	CassandraUsername string `mapstructure:"CASSANDRA_USERNAME" json:"CASSANDRA_USERNAME"`
	CassandraPassword string `mapstructure:"CASSANDRA_PASSWORD" json:"CASSANDRA_PASSWORD"`
}

func main() {
	app := cli.NewApp()
	app.Name = "Cassandra-Web"
	app.Version = "1.0.5"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Ken",
			Email: "ipushc@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "config.yaml",
			Usage:  "app config",
			EnvVar: "CONFIG_PATH",
		},
	}
	app.Action = run

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

// run
func run(c *cli.Context) {
	viper.SetConfigFile(c.String("config"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	envTmp := &envStruct{}
	err = viper.Unmarshal(envTmp)

	if err != nil {
		panic(err)
	}

	env = *envTmp

	log.Info("Cofing 設定成功")

	cluster := gocql.NewCluster(env.CassandraHost)
	cluster.Port = env.CassandraPort
	cluster.Keyspace = SystemSchemaKey
	cluster.Timeout = 30 * time.Second
	cluster.ConnectTimeout = 30 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 20}
	cluster.NumConns = 10
	cluster.Consistency = gocql.One

	if env.CassandraUsername != "" && env.CassandraPassword != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: env.CassandraUsername,
			Password: env.CassandraPassword,
		}
	}

	session, err := cluster.CreateSession()

	defer session.Close()

	if err != nil {
		log.Fatal(err)
	}

	h := &Handler{Session: session}

	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())

	// 跨網域用
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:8083", "http://localhost:8084"},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))

	// 讀靜態檔(前端)
	e.Static("/static", "client/dist/static")
	e.File("/", "client/dist/index.html")

	e.POST("/query", h.Query)
	e.POST("/save", h.Save)
	e.POST("/delete", h.Delete)
	e.POST("/find", h.Find)
	e.POST("/import", h.Import)
	e.POST("/rowtoken", h.RowToken)
	e.POST("/truncate", h.Truncate)

	e.GET("/keyspace", h.KeySpace)
	e.GET("/table", h.Table)
	e.GET("/row", h.Row)
	e.GET("/describe", h.Describe)
	e.GET("/columns", h.Columns)
	e.GET("/export", h.Export)
	e.GET("/hostinfo", h.HostInfo)

	// Start server
	e.Logger.Fatal(e.Start(env.HostPort))
}

type Handler struct {
	Session *gocql.Session
}

// Query Query cql語法處理
func (h *Handler) Query(c echo.Context) error {
	req := struct {
		Query string `json:"query" form:"query" query:"query"`
	}{}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var rets []interface{}

	for _, q := range strings.Split(req.Query, ";") {
		rowData := make([]map[string]interface{}, 0)

		if q == "" {
			continue
		}

		iter := h.Session.Query(q).Iter()

		ret, err := iter.SliceMap()

		for _, k := range ret {
			row := make(map[string]interface{})
			for vi, ki := range k {
				row[vi] = ki
			}

			rowData = append(rowData, OutputTransformType(row))
		}
		if err != nil {
			rets = append(rets, err.Error())
		}
		rets = append(rets, rowData)
	}

	return c.JSON(http.StatusOK, rets)
}

// KeySpace 取的所有keypace處理
func (h *Handler) KeySpace(c echo.Context) error {
	iter := h.Session.Query(`SELECT keyspace_name FROM system_schema.keyspaces`).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i, v := range ret {
		tableIter := h.Session.Query(`SELECT table_name FROM system_schema.tables WHERE  keyspace_name = ?`, v["keyspace_name"]).Iter()
		tableRet, err := tableIter.SliceMap()

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		for i, _ := range tableRet {
			tableRet[i]["views"] = false
		}

		viewIter := h.Session.Query(`SELECT view_name as table_name FROM system_schema.views WHERE  keyspace_name = ?`, v["keyspace_name"]).Iter()
		viewRet, err := viewIter.SliceMap()

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		for i, _ := range viewRet {
			viewRet[i]["views"] = true
		}

		tableRet = append(tableRet, viewRet...)

		ret[i]["table"] = tableRet

		// 避免前端element table出現關鍵字bug
		if v["keyspace_name"] == "system_distributed" {
			ret[i]["keyspace_name"] = "system_distributed!"
		}
	}

	return c.JSON(http.StatusOK, ret)
}

// Table 取的keyspace的table處理
func (h *Handler) Table(c echo.Context) error {
	keyspace := c.QueryParam("keyspace")

	// 查詢	table
	iter := h.Session.Query(`SELECT keyspace_name, table_name, id FROM system_schema.tables WHERE  keyspace_name = ?`, keyspace).Iter()
	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i, _ := range ret {
		ret[i]["views"] = false
	}

	// 查詢 虛擬view
	iter2 := h.Session.Query(`SELECT keyspace_name, view_name as table_name, id FROM system_schema.views WHERE  keyspace_name = ?`, keyspace).Iter()
	ret2, err := iter2.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i, _ := range ret2 {
		ret2[i]["views"] = true
	}

	ret = append(ret, ret2...)

	return c.JSON(http.StatusOK, ret)
}

type RowTokenReq struct {
	Table    string                 `json:"table" form:"table" query:"table"`
	Item     map[string]interface{} `json:"item" form:"item" query:"item"`
	PrevNext string                 `json:"prevnext" form:"prevnext" query:"prevnext"`
	Pagesize int                    `json:"pagesize" form:"pagesize" query:"pagesize"`
}

// HostInfo
func (h *Handler) HostInfo(c echo.Context) error {
	hosts, err := h.Session.GetHosts()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var ret []map[string]interface{}

	for _, host := range hosts {
		hostinfo := make(map[string]interface{})
		hostinfo["ClusterName"] = host.ClusterName()
		hostinfo["ConnectAddress"] = host.ConnectAddress()
		hostinfo["Peer"] = host.Peer()
		hostinfo["RPCAddress"] = host.RPCAddress()
		hostinfo["BroadcastAddress"] = host.BroadcastAddress()
		hostinfo["ListenAddress"] = host.ListenAddress()
		hostinfo["PreferredIP"] = host.PreferredIP()
		hostinfo["DataCenter"] = host.DataCenter()
		hostinfo["HostID"] = host.HostID()
		hostinfo["Port"] = host.Port()
		hostinfo["State"] = host.State().String()
		hostinfo["Version"] = host.Version().String()

		ret = append(ret, hostinfo)
	}

	return c.JSON(http.StatusOK, ret)
}

// RowToken
func (h *Handler) RowToken(c echo.Context) error {

	var (
		req    RowTokenReq
		schema []map[string]interface{}
	)

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	schema = h.GetSchema(req.Table)
	rowData, err := h.FirstQuery(&req, schema)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(rowData) < req.Pagesize && len(req.Item) >= 1 {
		sRowData, err := h.SecondQuery(&req, schema, req.Pagesize-len(rowData))

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		rowData = append(rowData, sRowData...)
	}

	data := make(map[string]interface{})
	data["row"] = rowData

	return c.JSON(http.StatusOK, data)
}
func (h *Handler) FirstQuery(req *RowTokenReq, schema []map[string]interface{}) ([]map[string]interface{}, error) {
	var (
		pCqlColumnName  []string
		pCqlColumnValue []interface{}
		pCqlPlaceholder []string
		cCqlColumnName  []string
		cCqlColumnValue []interface{}
		cCqlPlaceholder []string
		cPrevNext       string = ">"
		cql             string
		rowData         []map[string]interface{} = make([]map[string]interface{}, 0)
	)

	if len(req.Item) < 1 {
		cql = fmt.Sprintf("SELECT * FROM %s LIMIT %d", req.Table, req.Pagesize)
	} else {
		if req.PrevNext == "prev" {
			cPrevNext = "<"
		}

		for _, v := range schema {
			kind := v["kind"].(string)
			columnName := v["column_name"].(string)
			columnType := v["type"].(string)

			if _, ok := req.Item[columnName]; !ok {
				continue
			}

			if kind == PartitionKey {
				pCqlColumnName = append(pCqlColumnName, columnName)
				pCqlColumnValue = append(pCqlColumnValue, cqlFormatValue(columnType, req.Item[columnName]))
				pCqlPlaceholder = append(pCqlPlaceholder, "?")
			} else if kind == ClusteringKey {
				cCqlColumnName = append(cCqlColumnName, columnName)
				cCqlColumnValue = append(cCqlColumnValue, cqlFormatValue(columnType, req.Item[columnName]))
				cCqlPlaceholder = append(cCqlPlaceholder, "?")
			}
		}

		if len(cCqlColumnName) == 0 {
			return rowData, nil
		}

		cql = fmt.Sprintf("SELECT * FROM %s WHERE token(%s) %s token(%s) ", req.Table, strings.Join(pCqlColumnName, ","), "=", strings.Join(pCqlPlaceholder, ","))
		cql += fmt.Sprintf("AND (%s) %s (%s) LIMIT %d", strings.Join(cCqlColumnName, ","), cPrevNext, strings.Join(cCqlPlaceholder, ","), req.Pagesize)
	}
	fmt.Println(cql)
	rowIter := h.Session.Query(cql, append(pCqlColumnValue, cCqlColumnValue...)...).Iter()

	for {
		row := make(map[string]interface{})
		if !rowIter.MapScan(row) {
			break
		}
		rowData = append(rowData, OutputTransformType(row))
	}

	return rowData, nil
}

func (h *Handler) SecondQuery(req *RowTokenReq, schema []map[string]interface{}, limit int) ([]map[string]interface{}, error) {
	var (
		pCqlColumnName  []string
		pCqlColumnValue []interface{}
		pCqlPlaceholder []string
		pPrevNext       string = ">"
		cql             string
		rowData         []map[string]interface{} = make([]map[string]interface{}, 0)
	)

	if len(req.Item) < 1 {
		return rowData, nil
	} else {
		if req.PrevNext == "prev" {
			pPrevNext = "<"
		}

		for _, v := range schema {
			kind := v["kind"].(string)
			columnName := v["column_name"].(string)
			columnType := v["type"].(string)

			if _, ok := req.Item[columnName]; !ok {
				continue
			}

			if kind == PartitionKey {
				pCqlColumnName = append(pCqlColumnName, columnName)
				pCqlColumnValue = append(pCqlColumnValue, cqlFormatValue(columnType, req.Item[columnName]))
				pCqlPlaceholder = append(pCqlPlaceholder, "?")
			}
		}

		cql = fmt.Sprintf("SELECT * FROM %s WHERE token(%s) %s token(%s) ", req.Table, strings.Join(pCqlColumnName, ","), pPrevNext, strings.Join(pCqlPlaceholder, ","))
		cql += fmt.Sprintf("LIMIT %d", limit)
	}

	rowIter := h.Session.Query(cql, pCqlColumnValue...).Iter()

	for {
		row := make(map[string]interface{})
		if !rowIter.MapScan(row) {
			break
		}
		rowData = append(rowData, OutputTransformType(row))
	}

	return rowData, nil
}

// Row 取的table的row資料處理 (資量大時需耗費很多效能)
func (h *Handler) Row(c echo.Context) error {
	data := make(map[string]interface{})

	table := c.QueryParam("table")
	page, err := strconv.Atoi(c.QueryParam("page"))
	pagesize, err := strconv.Atoi(c.QueryParam("pagesize"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	countIter := h.Session.Query(`SELECT COUNT(*) FROM ` + table).Iter()
	countRet, err := countIter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data["count"] = countRet[0]["count"]

	if page == 0 {
		page = 1
	}

	limit_end := page * pagesize
	limit_start := limit_end - pagesize
	i := 0

	rowIter := h.Session.Query(`SELECT * FROM `+table+` LIMIT ?`, limit_end).Iter()
	rowData := make([]map[string]interface{}, 0)

	for {
		i++

		row := make(map[string]interface{})
		if !rowIter.MapScan(row) {
			break
		}
		if i > limit_start {
			rowData = append(rowData, OutputTransformType(row))
		}
	}

	data["row"] = rowData

	return c.JSON(http.StatusOK, data)
}

// Describe 用cqlsh取的describe
func (h *Handler) Describe(c echo.Context) error {
	table := c.QueryParam("table")

	cql := fmt.Sprintf("DESCRIBE %s ;", table)
	cmd := exec.Command("cqlsh", env.CassandraHost, "-e", cql)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, string(out))
}

// Columns 取的tabel欄位處理
func (h *Handler) Columns(c echo.Context) error {
	keyspace := c.QueryParam("keyspace")
	table := c.QueryParam("table")
	ret := h.GetSchema(keyspace + "." + table)

	return c.JSON(http.StatusOK, ret)
}

func (h *Handler) Save(c echo.Context) error {
	req := struct {
		Table string `json:"table" form:"table" query:"table"`
		Item  string `json:"item" form:"item" query:"item"`
	}{}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var item map[string]interface{}

	err := jsoni.Unmarshal([]byte(req.Item), &item)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var (
		itemKey         []string
		itemData        []interface{}
		itemPlaceholder []string
	)
	schema := make(map[string]string)

	for _, v := range h.GetSchema(req.Table) {
		schema[v["column_name"].(string)] = v["type"].(string)
	}

	itemKey, itemData, itemPlaceholder, err = InputTransformType(item, schema)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := h.Session.Query(`INSERT INTO `+req.Table+` (`+strings.Join(itemKey, ",")+`) VALUES(`+strings.Join(itemPlaceholder, ",")+`)`, itemData...).Exec(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// Delete 刪除row
func (h *Handler) Delete(c echo.Context) error {
	req := struct {
		Table string `json:"table" form:"table" query:"table"`
		Item  string `json:"item" form:"item" query:"item"`
	}{}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var item map[string]interface{}

	err := jsoni.Unmarshal([]byte(req.Item), &item)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	schema := h.GetSchema(req.Table)

	var (
		partitionCql    []string
		clusteringCql   []string
		partitionValue  []interface{}
		clusteringValue []interface{}
	)

	for _, v := range schema {
		kind := v["kind"].(string)
		columnName := v["column_name"].(string)
		columnType := v["type"].(string)

		if kind == PartitionKey {
			partitionCql = append(partitionCql, cqlFormatWhere(columnName, "="))
			partitionValue = append(partitionValue, cqlFormatValue(columnType, item[columnName]))
		} else if kind == ClusteringKey {
			clusteringCql = append(clusteringCql, cqlFormatWhere(columnName, "="))
			clusteringValue = append(clusteringValue, cqlFormatValue(columnType, item[columnName]))
		}
	}

	cql := fmt.Sprintf("DELETE FROM %s WHERE %s ", req.Table, strings.Join(append(partitionCql, clusteringCql...), " AND "))
	if err := h.Session.Query(cql, append(partitionValue, clusteringValue...)...).Exec(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// Find 搜尋row
func (h *Handler) Find(c echo.Context) error {
	req := struct {
		Table         string                            `json:"table" form:"table" query:"table"`
		Item          map[string]map[string]interface{} `json:"item" form:"item" query:"item"`
		OrderBy       []map[string]string               `json:"order_by" form:"order_by" query:"order_by"`
		Pagesize      int                               `json:"pagesize" form:"pagesize" query:"pagesize"`
		Page          int                               `json:"page" form:"page" query:"page"`
		IsAllowFilter bool                              `json:"isallowfilter" form:"isallowfilter" query:"isallowfilter"`
	}{}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	schema := h.GetSchema(req.Table)

	var (
		partitionCql    []string
		clusteringCql   []string
		partitionValue  []interface{}
		clusteringValue []interface{}
		orderByCql      []string
		AllowFilter     string
	)

	if req.IsAllowFilter {
		AllowFilter = "ALLOW FILTERING"
	}

	schemaMap := make(map[string]map[string]interface{})

	for _, v := range schema {
		schemaMap[v["column_name"].(string)] = v
	}

	for _, ob := range req.OrderBy {
		if _, ok := schemaMap[ob["name"]]; !ok {
			continue
		}

		if schemaMap[ob["name"]]["clustering_order"].(string) != "NONE" {
			orderByCql = append(orderByCql, fmt.Sprintf("%s %s", ob["name"], ob["order"]))
		}
	}

	if len(req.Item) == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, "not find condition")
	}

	for _, v := range schema {
		kind := v["kind"].(string)
		columnName := v["column_name"].(string)
		columnType := v["type"].(string)

		if kind == PartitionKey {
			if _, ok := req.Item[columnName]; !ok {
				continue
			}

			partitionCql = append(partitionCql, cqlFormatWhere(columnName, "="))
			partitionValue = append(partitionValue, cqlFormatValue(columnType, req.Item[columnName]["value"]))

		} else if kind == ClusteringKey || req.IsAllowFilter {
			if _, ok := req.Item[columnName]; !ok {
				continue
			}

			clusteringCql = append(clusteringCql, cqlFormatWhere(columnName, req.Item[columnName]["operator"].(string)))
			clusteringValue = append(clusteringValue, cqlFormatValue(columnType, req.Item[columnName]["value"]))
		}
	}

	conutCql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s %s", req.Table, strings.Join(append(partitionCql, clusteringCql...), " AND "), AllowFilter)
	countIter := h.Session.Query(conutCql, append(partitionValue, clusteringValue...)...).Iter()
	countRet, err := countIter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data := make(map[string]interface{})
	data["count"] = countRet[0]["count"]

	if req.Page == 0 {
		req.Page = 1
	}

	limit_end := req.Page * req.Pagesize
	limit_start := limit_end - req.Pagesize
	i := 0

	rowCql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT %d %s", req.Table, strings.Join(append(partitionCql, clusteringCql...), " AND "), limit_end, AllowFilter)
	rowIter := h.Session.Query(rowCql, append(partitionValue, clusteringValue...)...).Iter()
	rowData := make([]map[string]interface{}, 0)

	for {
		i++

		row := make(map[string]interface{})
		if !rowIter.MapScan(row) {
			break
		}
		if i > limit_start {
			rowData = append(rowData, OutputTransformType(row))
		}
	}

	data["row"] = rowData

	return c.JSON(http.StatusOK, data)
}

// Export 匯出copy file
func (h *Handler) Export(c echo.Context) error {
	table := c.QueryParam("table")

	cql := fmt.Sprintf("COPY %s TO STDOUT;", table)
	cmd := exec.Command("cqlsh", env.CassandraHost, "-e", cql)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", "cql", strings.Replace(table, ".", "-", -1)))
	return c.Blob(http.StatusOK, "application/force-download", out)
}

// Export 匯入copy file
func (h *Handler) Import(c echo.Context) error {
	file, err := c.FormFile("file")
	table := c.FormValue("table")

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	f, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer f.Close()

	tmpPath := "/tmp/importfile"
	if _, err := os.Stat(tmpPath); err != nil {
		if os.IsNotExist(err) {
			CreateTmpFile(tmpPath)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	tf, err := os.OpenFile(
		tmpPath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tf.Close()

	fb, err := ioutil.ReadAll(f)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = tf.Write(fb)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cql := fmt.Sprintf("COPY %s FROM '%s' WITH MINBATCHSIZE=1 AND MAXBATCHSIZE=1 AND PAGESIZE=10;", table, tmpPath)
	cmd := exec.Command("cqlsh", env.CassandraHost, "--connect-timeout=600", "--request-timeout=600", "-e", cql)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// Truncate 清除table資料
func (h *Handler) Truncate(c echo.Context) error {
	req := struct {
		Table string `json:"table" form:"table" query:"table"`
	}{}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cql := fmt.Sprintf("TRUNCATE TABLE %s;", req.Table)

	if err := h.Session.Query(cql).Exec(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// GetSchema 取的table schema
func (h *Handler) GetSchema(table string) []map[string]interface{} {
	tablekey := strings.Split(table, ".")

	iter := h.Session.Query(`SELECT * FROM system_schema.columns WHERE keyspace_name = ? and table_name = ? `, tablekey[0], tablekey[1]).Iter()
	ret, err := iter.SliceMap()

	if err != nil {
		return nil
	}

	sort.Slice(ret, func(a, b int) bool {
		if ret[a]["kind"].(string) == PartitionKey {
			if ret[b]["kind"].(string) == PartitionKey {
				return ret[a]["position"].(int) < ret[b]["position"].(int)
			}

			return true
		}

		if ret[a]["kind"].(string) == ClusteringKey {
			if ret[b]["kind"].(string) == PartitionKey {
				return false
			} else if ret[b]["kind"].(string) == ClusteringKey {
				return ret[a]["position"].(int) < ret[b]["position"].(int)
			}

			return true
		}

		return false
	})

	return ret
}
