package handlers

import (
	"httpserver/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	db *database.Database
}

func (s *server) GetUserSegments(c *gin.Context) {

	ctx := c.Request.Context()

	userID := c.Param("id")

	segments, err := s.db.GetUserSegments(ctx, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, segments)
}
