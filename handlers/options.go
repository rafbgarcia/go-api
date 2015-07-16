package handlers

import (
	"github.com/gin-gonic/gin"
)

var allowedDomains = []string{
	"http://localhost:8080",
	"https://go-web-rafa.appspot.com",
}

func SetOrigin(c *gin.Context) {
	for _, domain := range allowedDomains {
		if domain == c.Request.Header.Get("Origin") {
			c.Header("Access-Control-Allow-Origin", domain)
			return
		}
	}
}

func OptionsHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "PUT, DELETE")
}
