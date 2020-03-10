package main

import (
	"database/sql"
	"fmt"
	"github.com/appleboy/go-fcm"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	// debug
	const (
		host     = "localhost"
		port     = 5432
		user     = "nomeet"
		password = "nomeet"
		dbname   = "coronavirus"
	)

	// release
	//const (
	//	host     = "localhost"
	//	port     = 5432
	//	user     = "postgres"
	//	password = "postgres"
	//	dbname   = "coronavirus"
	//)

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

	client, err := fcm.NewClient("AAAAY5aoA3U:APA91bGo_vKVUguG0wIqysIZq7kHC1Fu1AbwzNJzS2uv0TkACLGAFKTb8USCFX8bS91U9HRM06TC7eylS8hFbTdxVT1zZQVwVn2FTCU3IunD4YZ9wpOLQa0eXbVzizXGQv5-LZrh7F-I")
	if err != nil {
		log.Fatalln(err)
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
				//"date" : date.Format("02-Jan-2006 15:04:05"),
				"date" : date.Format("Jan 02, 15:04 EST") + date.Location().String(),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success" : true,
			"results" : results,
		})
	})

	r.POST("/sendNotif", func(c *gin.Context) {
		if c.PostForm("password") != "juanisdumb" {
			c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("wrong password"))
			return
		}
		msg := &fcm.Message{
			Notification: &fcm.Notification{
				Title: c.PostForm("title"),
				Body: c.PostForm("description"),
			},
			Condition: "'" + c.PostForm("topic") + "' in Topics",
			Data: map[string]interface{}{
				"message": "yes",
			},
		}

		// Send the message and receive the response without retries.
		response, err := client.Send(msg)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%#v\n", response)

		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte("success"))
		return
	})

	r.RunTLS(":443", "key/domain-crt.txt", "key/account-key.txt")
}