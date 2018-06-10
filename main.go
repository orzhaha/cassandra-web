package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

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
	// 	AllowOrigins: []string{"http://localhost:8081", "http://localhost:8082"},
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

	iter := h.Session.Query(`SELECT * FROM system_schema.tables WHERE  keyspace_name = ?`, keyspace).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

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
			rowData = append(rowData, row)
		}
	}

	data["row"] = rowData

	return c.JSON(http.StatusOK, data)
}

type SaveReq struct {
	Item  string `json:"item" form:"item" query:"item"`
	Table string `json:"table" form:"table" query:"table"`
}

func (h *Handler) Save(c echo.Context) error {
	saveReq := new(SaveReq)

	if err := c.Bind(saveReq); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var item map[string]interface{}

	err := json.Unmarshal([]byte(saveReq.Item), &item)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var (
		itemKey         []string
		itemData        []interface{}
		itemPlaceholder []string
	)

	for k, v := range item {
		itemKey = append(itemKey, k)
		itemData = append(itemData, v)
		itemPlaceholder = append(itemPlaceholder, "?")
	}

	if err := h.Session.Query(`INSERT INTO `+saveReq.Table+` (`+strings.Join(itemKey, ",")+`) VALUES(`+strings.Join(itemPlaceholder, ",")+`)`, itemData...).Exec(); err != nil {
		log.Info(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}
