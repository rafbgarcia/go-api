package posts

import (
	"net/http"

	"appengine"
	"appengine/datastore"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Id    string `json:"id" datastore:"-"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func Delete(c *gin.Context) {
	gaeContext := appengine.NewContext(c.Request)

	key, _ := datastore.DecodeKey(c.Param("id"))

	if err := datastore.Delete(gaeContext, key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func Create(c *gin.Context) {
	gaeContext := appengine.NewContext(c.Request)
	post := new(Post)

	c.Bind(&post)

	key := datastore.NewIncompleteKey(gaeContext, "Posts", nil)

	if _, err := datastore.Put(gaeContext, key, post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func Edit(c *gin.Context) {
	gaeContext := appengine.NewContext(c.Request)
	post := new(Post)

	key, _ := datastore.DecodeKey(c.Param("id"))

	if err := datastore.Get(gaeContext, key, post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Bind(&post)

	if _, err := datastore.Put(gaeContext, key, post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func List(c *gin.Context) {
	gaeContext := appengine.NewContext(c.Request)
	posts := []*Post{}

	keys, err := datastore.NewQuery("Posts").GetAll(gaeContext, &posts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i, post := range posts {
		post.Id = keys[i].Encode()
	}

	c.JSON(http.StatusOK, posts)
}
