package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"httpserver/database"
	"httpserver/handler"
	"httpserver/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/jackc/pgx/v5"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func router() *gin.Engine {
	router := gin.Default()
	h, _ := handler.New(db)
	router.GET("/user/:id", h.GetUserSegments)
	router.PATCH("/user/:id", h.UpdateUserSegments)
	router.POST("/segment/create", h.CreateSegment)
	router.DELETE("/segment/:id", h.DeleteSegment)

	return router
}

var db *database.Database

func setup() {
	dsn := "postgres://postgres:postgres@localhost:5432/postgres_test?sslmode=disable"
	db, _ = database.New(dsn)

	ctx := context.Background()
	tx, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	query := `

CREATE TABLE IF NOT EXISTS users(
                                    user_id SERIAL PRIMARY KEY,
                                    username VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS segments(
                                       seg_id SERIAL PRIMARY KEY,
                                       seg_name VARCHAR(50) --TODO: сделать текстом
);

CREATE TABLE IF NOT EXISTS user_segment (
                                            user_id INT,
                                            segment_id INT,
                                            PRIMARY KEY (user_id, segment_id),
                                            FOREIGN KEY (user_id) REFERENCES users(user_id),
                                            FOREIGN KEY (segment_id) REFERENCES segments(seg_id)
);

	`
	_, err = tx.Exec(ctx, query)
	if err != nil {
		_ = tx.Rollback(ctx)
		return
	}
	tx.Commit(ctx)

}

func teardown() {
	defer db.Close(context.Background())

	ctx := context.Background()
	tx, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	query := `
    DROP TABLE user_segment cascade;
    DROP TABLE users cascade;
    DROP TABLE segments cascade;
	`
	_, err = tx.Exec(ctx, query)
	if err != nil {
		fmt.Println(err.Error())
		_ = tx.Rollback(ctx)
		return
	}
	tx.Commit(ctx)

}

func makeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

func TestCreateSegment(t *testing.T) {
	seg := model.SegName{
		SegName: "AVITO_VOICE_MESSAGES",
	}
	writer := makeRequest("POST", "/segment/create", seg)
	assert.Equal(t, http.StatusOK, writer.Code)

	seg1 := model.SegName{
		SegName: "AVITO_PERFORMANCE_VAS",
	}
	writer1 := makeRequest("POST", "/segment/create", seg1)
	assert.Equal(t, http.StatusOK, writer1.Code)
}

func TestCreateSegmentExist(t *testing.T) {
	seg := model.SegName{
		SegName: "AVITO_VOICE_MESSAGES",
	}
	writer := makeRequest("POST", "/segment/create", seg)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestDeleteSegment(t *testing.T) {
	seg := model.SegName{
		SegName: "AVITO_VOICE_MESSAGES",
	}
	writer := makeRequest("DELETE", "/segment/1", seg)
	assert.Equal(t, http.StatusOK, writer.Code)
}