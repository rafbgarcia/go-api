package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "appengine/datastore"
    "appengine"
)

func init() {
    router := gin.Default()

    router.GET("/posts", listPosts)
    router.POST("/posts", createPost)

    http.Handle("/", router)
}


type Post struct {
    Title string `json:"title"`
    Body string `json:"body"`
}

func createPost(c *gin.Context) {
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

func listPosts(c *gin.Context) {
    gaeContext := appengine.NewContext(c.Request)
    posts := []*Post{}

    if _, err := datastore.NewQuery("Posts").GetAll(gaeContext, &posts); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, posts)
}