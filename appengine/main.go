package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"appengine"
	"appengine/datastore"
	"crypto/rand"
	"fmt"
	"github.com/rafbgarcia/go-api/handlers"
	"github.com/rafbgarcia/go-api/posts"
)

func init() {
	router := gin.Default()

	router.Use(handlers.SetOrigin)

	router.OPTIONS("/posts", handlers.OptionsHandler)
	router.OPTIONS("/posts/:id", handlers.OptionsHandler)

	router.GET("/posts", posts.List)

	{
		authRouter := router.Group("/posts", authenticate)

		authRouter.POST("", posts.Create)
		authRouter.PUT("/:id", posts.Edit)
		authRouter.DELETE("/:id", posts.Delete)
	}

	router.POST("users", createUser)

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

	// username:password
	_, apiKey, ok := c.Request.BasicAuth()

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid API Key",
		})
		c.Abort()
		return
	}

	user := new(User)
	key := datastore.NewKey(gaeContext, "Users", apiKey, 0, nil)

	if err := datastore.Get(gaeContext, key, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid API Key",
		})
		c.Abort()
		return
	}
}
