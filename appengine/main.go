package main

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/rafbgarcia/go-api/posts"
    "github.com/rafbgarcia/go-api/handlers"

    "appengine/datastore"
    "appengine"
    "crypto/rand"
    "fmt"
)

func init() {
    router := gin.Default()

    router.Use(handlers.SetOrigin)
    router.Use(handlers.SetOrigin)

    router.OPTIONS("/posts", handlers.OptionsHandler)

    router.GET("/posts", posts.List)
    router.POST("/users", createUser)

    admin := router.Group("/posts", authenticate)
    {
        admin.POST("", posts.Create)
        admin.PUT("/:id", posts.Edit)
        admin.DELETE("/:id", posts.Delete)
    }

    http.Handle("/", router)
}

type User struct {
    ApiKey string `json:"api_key"`
}

func createUser(c *gin.Context) {
    gaeContext := appengine.NewContext(c.Request)

    bytes := make([]byte, 40)
    rand.Read(bytes)

    user := new(User)
    user.ApiKey = fmt.Sprintf("%x", bytes)

    key := datastore.NewKey(gaeContext, "Users", user.ApiKey, 0, nil)

    if _, err := datastore.Put(gaeContext, key, user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, user)
}

func authenticate(c *gin.Context) {
    gaeContext := appengine.NewContext(c.Request)

    _, apiKey, ok := c.Request.BasicAuth()

    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid API Key",
        })
        c.Abort()
        return
    }

    key := datastore.NewKey(gaeContext, "Users", apiKey, 0, nil)

    user := new(User)

    if err := datastore.Get(gaeContext, key, user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid API Key",
        })
        c.Abort()

        return
    }
}
