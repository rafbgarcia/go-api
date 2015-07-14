package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "appengine/datastore"
    "appengine"
)

func init() {
    router := gin.Default()

    router.Use(setOrigin)

    router.OPTIONS("/posts", optionsHandler)
    router.OPTIONS("/posts/:id", optionsHandler)

    router.GET("/posts", listPosts)
    router.POST("/posts", createPost)
    router.PUT("/posts/:id", editPost)
    router.DELETE("/posts/:id", deletePost)

    http.Handle("/", router)
}


type Post struct {
    Id string `json:"id" datastore:"-"`
    Title string `json:"title"`
    Body string `json:"body"`
}

func deletePost(c *gin.Context) {
    gaeContext := appengine.NewContext(c.Request)

    key, _ := datastore.DecodeKey(c.Param("id"))

    if err := datastore.Delete(gaeContext, key); err !=nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.AbortWithStatus(http.StatusNoContent)
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

func editPost(c *gin.Context) {
    gaeContext := appengine.NewContext(c.Request)
    post := new(Post)

    key, _ := datastore.DecodeKey(c.Param("id"))

    if err := datastore.Get(gaeContext, key, post); err !=nil {
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

func listPosts(c *gin.Context) {
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

var allowedDomains = []string{
    "http://localhost:8080",
    "https://go-web-rafa.appspot.com",
}

func setOrigin(c *gin.Context) {
    for _, domain := range allowedDomains {
        if domain == c.Request.Header.Get("Origin") {
            c.Header("Access-Control-Allow-Origin", domain)
            break
        }
    }
}

func optionsHandler(c *gin.Context) {
    c.Header("Access-Control-Allow-Headers", "Content-Type")
    c.Header("Access-Control-Allow-Methods", "PUT, DELETE")
}