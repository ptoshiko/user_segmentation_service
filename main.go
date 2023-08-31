package main

import (
	"context"
	"regexp"

	"httpserver/database"
	"httpserver/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type server struct {
	db *database.Database
}

var validSegmentNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func isValidSegmentName(name string) bool {
	return validSegmentNamePattern.MatchString(name) && name != ""
}

func (s *server) createSegment(c *gin.Context) {

	ctx := c.Request.Context()

	var seg model.SegName

	if err := c.ShouldBindJSON(&seg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if !isValidSegmentName(seg.SegName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment name"})
		return
	}

	err := s.db.CreateSegment(ctx, seg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Segment created successfully"})
}

// Метод удаления сегмента. Принимает slug (название) сегмента.

func (s *server) deleteSegment(c *gin.Context) {
	ctx := c.Request.Context()

	var seg model.SegName
	if err := c.ShouldBindJSON(&seg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := s.db.DeleteSegment(ctx, seg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Segment deleted successfully"})
}

// Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю,
// список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.

func (s *server) updateUserSegments(c *gin.Context) {
	ctx := c.Request.Context()

	var us model.UserSegments
	if err := c.ShouldBindJSON(&us); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := s.db.UpdateUserSegments(ctx, us)
	if err != nil {
		log.Printf("db.UpdateUserSegments: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User segments updated successfully"})
}

// Метод получения активных сегментов пользователя. Принимает на вход id пользователя.

func (s *server) getUserSegments(c *gin.Context) {

	ctx := c.Request.Context()

	userID := c.Param("id")

	segments, err := s.db.GetUserSegments(ctx, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, segments)
}

func main() {

	db, err := database.New()
	defer db.Close(context.Background())
	if err != nil {
		return
	}
	s := server{db: db}

	// dsn := "postgres://postgres:postgres@postgres:5432/postgres" + "?sslmode=disable"

	// conn, err := pgx.Connect(context.Background(), dsn)
	// if err != nil {
	// 	log.Fatalf("Unable to connect to database: %v\n", err)
	// }
	// defer conn.Close(context.Background())

	// if err = conn.Ping(context.Background()); err != nil {
	// 	log.Fatalf("can't ping db: %s", err)
	// }

	// s := server{conn: conn}

	router := gin.Default()

	// router := gin.New()
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	router.GET("/user/:id", s.getUserSegments)
	router.PATCH("/user/:id", s.updateUserSegments)
	router.POST("/segment/create", s.createSegment)
	router.DELETE("/segment/:id", s.deleteSegment)

	log.Fatal(router.Run(":8080"))
}