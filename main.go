package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-contrib/static"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(gin.Logger())
	r.Delims("{{", "}}")

	r.LoadHTMLGlob("./templates/*.tmpl.html")

	r.Use(static.Serve("/assets", static.LocalFile("./assets", false)))

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

		//p := bluemonday.UGCPolicy()
		//p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
		//html := p.SanitizeBytes(unsafe)

		p := bluemonday.UGCPolicy()
		p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")

		unsafe := blackfriday.Run([]byte(mdfile))
		//unsafe := blackfriday.MarkdownCommon([]byte(mdfile))
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

		postHTML := template.HTML(html)

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
