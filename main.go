package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"time"
)

func main() {

	const (
		host     = "localhost"
		port     = 5432
		user     = "nomeet"
		password = "nomeet"
		dbname   = "coronavirus"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("connected to postgres")

	r := gin.Default()

	r.Use(static.Serve("/secret/", static.LocalFile("./secret", true)))
	r.Use(static.Serve("/static/", static.LocalFile("./static", true)))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/add", func(c *gin.Context) {
		pwd := c.PostForm("password")
		if pwd == "mkesjddn" {
			title := c.PostForm("title")
			url := c.PostForm("url")
			_, err := db.Exec("INSERT INTO news (title, url) VALUES ($1, $2);", title, url)
			if err != nil {
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("error"))
			} else {
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("success"))
			}
		}
	})

	r.GET("/latest/:offset", func(c *gin.Context) {

		const itemLimit = 10
		offset := 0
		var err error
		if c.Param("offset") != "" {
			offset,  err = strconv.Atoi(c.Param("offset"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success" : false,
					"error" : "offset must be a number",
				})
				return
			}
		}
		fmt.Println("offset: ")
		fmt.Println(offset)
		query := "SELECT * FROM news LIMIT $1 OFFSET $2;"
		rows, err := db.Query(query, itemLimit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success" : false,
				"error" : "Database error. Anna oop 3.",
			})
			fmt.Println(err.Error())
			return
		}


		var id, title, url string
		var date time.Time
		var results [itemLimit]gin.H

		for i := 0; rows.Next(); i++ {
			rows.Scan(&id, &title, &url, &date)
			results[i] = gin.H{
				"id" : id,
				"title" : title,
				"url" : url,
				"date" : date.Format("02-Jan-2006 15:04:05"),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success" : true,
			"results" : results,
		})
	})

	r.Run()
}