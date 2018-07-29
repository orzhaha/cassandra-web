package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/gocql/gocql"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func init() {
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

	e.Static("/static", "client/dist/static")
	e.File("/", "client/dist/index.html")

	e.POST("/query", h.Query)
	e.POST("/save", h.Save)

	e.GET("/keyspace", h.KeySpace)
	e.GET("/table", h.Table)
	e.GET("/row", h.Row)

	// Start server
	e.Logger.Fatal(e.Start(env.HostPort))
}

type Handler struct {
	Session *gocql.Session
}

type CqlQuery struct {
	Query string `json:"query" form:"query" query:"query"`
}

func (h *Handler) Query(c echo.Context) error {
	query := new(CqlQuery)

	if err := c.Bind(query); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var rets []interface{}

	for _, q := range strings.Split(query.Query, ";") {
		if q == "" {
			continue
		}

		iter := h.Session.Query(q).Iter()

		ret, err := iter.SliceMap()

		if err != nil {
			rets = append(rets, err.Error())
		}

		rets = append(rets, ret)
	}

	return c.JSON(http.StatusOK, rets)
}

func (h *Handler) KeySpace(c echo.Context) error {
	iter := h.Session.Query(`SELECT * FROM system_schema.keyspaces`).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ret)
}

func (h *Handler) Table(c echo.Context) error {
	keyspace := c.QueryParam("keyspace")

	iter := h.Session.Query(`SELECT keyspace_name, table_name, id FROM system_schema.tables WHERE  keyspace_name = ?`, keyspace).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	iter2 := h.Session.Query(`SELECT keyspace_name,  view_name as table_name, id FROM system_schema.views WHERE  keyspace_name = ?`, keyspace).Iter()
	ret2, err := iter2.SliceMap()

	ret = append(ret, ret2...)

	return c.JSON(http.StatusOK, ret)
}

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

	rowIter := h.Session.Query(`SELECT * FROM `+table+` LIMIT ? ALLOW FILTERING`, limit_end).Iter()
	rowData := make([]map[string]interface{}, 0)

	for {
		i++

		row := make(map[string]interface{})
		if !rowIter.MapScan(row) {
			break
		}
		if i > limit_start {
			rowData = append(rowData, OutTransformType(row))
		}
	}

	data["row"] = rowData

	return c.JSON(http.StatusOK, data)
}

func OutTransformType(row map[string]interface{}) map[string]interface{} {
	for k, v := range row {
		switch v.(type) {
		case int64, float64, float32:
			row[k] = cast.ToString(v)
		case []int64:
			row[k] = cast.ToStringSlice(v)
		case map[string]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int64]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int32]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int16]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int8]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[float64]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[float32]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[bool]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		}

	}

	return row
}

type SaveReq struct {
	Table string `json:"table" form:"table" query:"table"`
	Item  string `json:"item" form:"item" query:"item"`
}

func (h *Handler) Save(c echo.Context) error {
	req := new(SaveReq)

	if err := c.Bind(req); err != nil {
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

	itemKey, itemData, itemPlaceholder, err = InTransformType(item, schema)

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

func InTransformType(item map[string]interface{}, schema map[string]string) ([]string, []interface{}, []string, error) {
	var (
		itemKey         []string
		itemData        []interface{}
		itemPlaceholder []string
	)

	mapReg := regexp.MustCompile(`(?U)^map\<(.+),\s(.+)\>`)

	for k, v := range item {
		switch schema[k] {
		case "bigint":
			val, err := cast.ToInt64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "boolean":
			val, err := cast.ToBoolE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "counter":
			val, err := cast.ToInt64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "date":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "decimal":
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "double":
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "float":
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "inet":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "int":
			val, err := cast.ToInt32E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "smallint":
			val, err := cast.ToInt16E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "text":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "time":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "timestamp":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "timeuuid":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "tinyint":
			val, err := cast.ToInt8E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "uuid":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "varchar":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "varint":
			val, err := cast.ToInt32E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		default:
			var val interface{} = v
			var err error

			mapRet := mapReg.FindStringSubmatch(schema[k])

			if len(mapRet) == 3 {
				val, err = MapToCassandraMapType(v, mapRet[1], mapRet[2])

				if err != nil {
					return nil, nil, nil, err
				}
			}

			itemData = append(itemData, val)
		}

		itemKey = append(itemKey, k)
		itemPlaceholder = append(itemPlaceholder, "?")
	}

	return itemKey, itemData, itemPlaceholder, nil
}

// MapToCassandraMapType map對應cassandra map的型別
func MapToCassandraMapType(i interface{}, keyType string, valType string) (interface{}, error) {
	var m = map[interface{}]interface{}{}

	switch v := i.(type) {
	case map[string]interface{}:
		for k, val := range v {
			kRet, err := CassandraTypeToGoType(k, keyType)

			if err != nil {
				return nil, err
			}

			valRet, err := CassandraTypeToGoType(val, valType)

			if err != nil {
				return nil, err
			}

			m[kRet] = valRet
		}

		return m, nil
	case string:
		err := JsonStringToObject(v, &m)
		return m, err
	case nil:
		return m, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]string", i, i)
	}
}

// CassandraTypeToGoType cassandra的型別轉Go型別
func CassandraTypeToGoType(i interface{}, t string) (interface{}, error) {
	switch t {
	case "bigint":
		val, err := cast.ToInt64E(i)
		return val, err
	case "boolean":
		val, err := cast.ToBoolE(i)

		return val, err
	case "counter":
		val, err := cast.ToInt64E(i)

		return val, err
	case "date":
		val, err := cast.ToStringE(i)

		return val, err
	case "decimal":
		val, err := cast.ToFloat64E(i)

		return val, err
	case "double":
		val, err := cast.ToFloat64E(i)

		return val, err
	case "float":
		val, err := cast.ToFloat64E(i)

		return val, err
	case "inet":
		val, err := cast.ToStringE(i)

		return val, err
	case "int":
		val, err := cast.ToInt32E(i)

		return val, err
	case "smallint":
		val, err := cast.ToInt16E(i)

		return val, err
	case "text":
		val, err := cast.ToStringE(i)

		return val, err
	case "time":
		val, err := cast.ToStringE(i)

		return val, err
	case "timestamp":
		val, err := cast.ToStringE(i)

		return val, err
	case "timeuuid":
		val, err := cast.ToStringE(i)

		return val, err
	case "tinyint":
		val, err := cast.ToInt8E(i)

		return val, err
	case "uuid":
		val, err := cast.ToStringE(i)

		return val, err
	case "varchar":
		val, err := cast.ToStringE(i)

		return val, err
	case "varint":
		val, err := cast.ToInt32E(i)

		return val, err
	}

	return i, nil
}

func JsonStringToObject(s string, v interface{}) error {
	data := []byte(s)
	return jsoni.Unmarshal(data, v)
}
