package main

import (
	"context"

	"httpserver/database"
	"httpserver/handler"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	dsn := "postgres://postgres:postgres@postgres:5432/postgres" + "?sslmode=disable"
	db, err := database.New(dsn)
	defer db.Close(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return
	}
	h, err := handler.New(db)
	router := gin.Default()

	router.GET("/user/:id", h.GetUserSegments)
	router.PATCH("/user/:id", h.UpdateUserSegments)
	router.POST("/segment/create", h.CreateSegment)
	router.DELETE("/segment/delete", h.DeleteSegment)

	log.Fatal(router.Run(":8080"))
}
