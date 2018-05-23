package main

import (
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	SystemSchemaKey = "system_schema"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())

	// 跨網域用
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:80"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Static("/static", "client/dist/static")
	e.File("/", "client/dist/index.html")

	e.POST("/query", postQuery)

	e.GET("/allkeyspace", getAllKeySpace)
	e.GET("/alltablebykeyspace", getAllTableByKeySpace)
	e.GET("/allrowbytable", getAllRowByTable)

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}

func postQuery(c echo.Context) error {
	cluster := gocql.NewCluster("cassandra")
	cluster.Port = 9042
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()

	defer session.Close()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	query := c.FormValue("query")

	iter := session.Query(query).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ret)
}

func getAllKeySpace(c echo.Context) error {
	cluster := gocql.NewCluster("cassandra")
	cluster.Port = 9042
	cluster.Keyspace = SystemSchemaKey
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()

	defer session.Close()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	iter := session.Query(`SELECT * FROM system_schema.keyspaces`).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ret)
}

func getAllTableByKeySpace(c echo.Context) error {
	cluster := gocql.NewCluster("cassandra")
	cluster.Port = 9042
	cluster.Keyspace = SystemSchemaKey
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()

	defer session.Close()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	keyspace := c.QueryParam("keyspace")

	iter := session.Query(`SELECT * FROM system_schema.tables WHERE  keyspace_name = ?`, keyspace).Iter()

	ret, err := iter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ret)
}

func getAllRowByTable(c echo.Context) error {
	cluster := gocql.NewCluster("cassandra")
	cluster.Port = 9042
	cluster.Keyspace = SystemSchemaKey
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	session.SetPageSize(200)

	defer session.Close()

	data := make(map[string]interface{})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	table := c.QueryParam("table")
	tokenKey := c.QueryParam("token_key")
	nextToken := c.QueryParam("next_token")
	prevToken := c.QueryParam("prev_token")
	limit, err := strconv.Atoi(c.QueryParam("limit"))

	countIter := session.Query(`SELECT COUNT(*) FROM ` + table).Iter()
	countRet, err := countIter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data["count"] = countRet[0]["count"]

	if nextToken != "" {
		rowIter := session.Query(`SELECT * FROM `+table+` WHERE token(`+tokenKey+`) > token('`+nextToken+`') LIMIT ? ALLOW FILTERING`, limit).Iter()
		rowRet, nextErr := rowIter.SliceMap()

		if nextErr != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, nextErr.Error())
		}

		data["row"] = rowRet

		return c.JSON(http.StatusOK, data)
	}

	if prevToken != "" {
		rowIter := session.Query(`SELECT * FROM `+table+` WHERE token(`+tokenKey+`) < token('`+prevToken+`') LIMIT ? ALLOW FILTERING`, limit).Iter()
		rowRet, prevErr := rowIter.SliceMap()

		if prevErr != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, prevErr.Error())
		}

		data["row"] = rowRet

		return c.JSON(http.StatusOK, data)
	}

	rowIter := session.Query(`SELECT * FROM `+table+` LIMIT ? ALLOW FILTERING`, limit).Iter()
	rowRet, err := rowIter.SliceMap()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data["row"] = rowRet

	return c.JSON(http.StatusOK, data)
}
