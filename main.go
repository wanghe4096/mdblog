package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(gin.Logger())
	r.Delims("{{", "}}")

	r.LoadHTMLGlob("./templates/*.tmpl.html")

	r.GET("/", func(c *gin.Context) {
		var posts []string

		files, err := ioutil.ReadDir("./_posts/")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
			posts = append(posts, file.Name())
		}

		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"posts": posts,
		})
	})

	r.GET("/:postName", func(c *gin.Context) {
		postName := c.Param("postName")

		mdfile, err := ioutil.ReadFile("./_posts/" + postName)

		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusNotFound, "error.tmpl.html", nil)
			c.Abort()
			return
		}

		postHTML := template.HTML(blackfriday.MarkdownCommon([]byte(mdfile)))

		post := Post{Title: postName, Content: postHTML}

		c.HTML(http.StatusOK, "post.tmpl.html", gin.H{
			"Title":   post.Title,
			"Content": post.Content,
		})
	})

	r.Run()
}

type Post struct {
	Title   string
	Content template.HTML
}
