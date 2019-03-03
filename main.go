package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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
}

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	SystemSchemaKey = "system_schema"
)

var env envStruct

// envStruct type
type envStruct struct {
	HostPort      string `mapstructure:"HOST_PORT" json:"HOST_PORT"`
	CassandraHost string `mapstructure:"CASSANDRA_HOST" json:"CASSANDRA_HOST"`
	CassandraPort int    `mapstructure:"CASSANDRA_PORT" json:"CASSANDRA_PORT"`
}

func main() {
	app := cli.NewApp()
	app.Name = "Cassandra-Web"
	app.Version = "1.0.2"
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

	e.GET("/keyspace", h.KeySpace)
	e.GET("/table", h.Table)
	e.GET("/row", h.Row)
	e.GET("/describe", h.Describe)
	e.GET("/columns", h.Columns)

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

	for i, v := range ret {
		// 避免前端element table出現關鍵字bug
		if v["keyspace_name"] == "system_distributed" {
			ret[i]["keyspace_name"] = "system_distributed!"
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

	// 查詢 虛擬view
	iter2 := h.Session.Query(`SELECT keyspace_name, view_name as table_name, id FROM system_schema.views WHERE  keyspace_name = ?`, keyspace).Iter()
	ret2, err := iter2.SliceMap()

	ret = append(ret, ret2...)

	return c.JSON(http.StatusOK, ret)
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
	kind := c.QueryParam("kind")
	item := c.QueryParam("item")

	cql := fmt.Sprintf("DESCRIBE %s %s ;", kind, item)
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

	cql := fmt.Sprintf("SELECT * FROM system_schema.columns WHERE keyspace_name='%s' AND table_name='%s';", keyspace, table)
	iter := h.Session.Query(cql).Iter()
	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

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

func (h *Handler) GetSchema(table string) []map[string]interface{} {
	tablekey := strings.Split(table, ".")

	iter := h.Session.Query(`SELECT * FROM system_schema.columns WHERE keyspace_name = '` + tablekey[0] + `' and table_name = '` + tablekey[1] + `' `).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return nil
	}

	return ret
}
