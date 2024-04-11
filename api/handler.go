package api

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) InitHandlers() {
	s.GET("/hello", s.hello)
}

func (s *Service) hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello",
	})
	return
}
