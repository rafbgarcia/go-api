package main

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/rafbgarcia/go-api/posts"
    "github.com/rafbgarcia/go-api/handlers"
)

func init() {
    router := gin.Default()

    router.Use(handlers.SetOrigin)

    router.OPTIONS("/posts", handlers.OptionsHandler)
    router.OPTIONS("/posts/:id", handlers.OptionsHandler)

    router.GET("/posts", posts.List)
//    router.POST("/posts", posts.Create)
//    router.PUT("/posts/:id", posts.Edit)
//    router.DELETE("/posts/:id", posts.Delete)

    http.Handle("/", router)
}
