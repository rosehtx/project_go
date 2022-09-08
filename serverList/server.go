package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"serverList/router"
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"`
}

func startServer() {
	r := router.InitRouter()
	
	//r.POST("/gin/test",test)
	r.Run(":8090")
}

func test(c *gin.Context) {
	var json Login
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(json.Password)
	if json.Password == "" {
		c.JSON(http.StatusOK, gin.H{"status": "empty password"})
		return
	}
	if json.User != "aa" || json.Password != "bb" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}
