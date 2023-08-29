package main

import (
	"context"
	"regexp"

	// "httpserver/database"
	"httpserver/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type server struct {
	conn *pgx.Conn
}

// createSegmnet moment
func responseJSON(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return
	}
	defer tx.Rollback(ctx)
	// Проверка, что такого сегмента еще нет
	var id int
	err = tx.QueryRow(ctx, `SELECT seg_id FROM segments WHERE seg_name = $1`, seg.SegName).Scan(&id)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Segment already exists", "id": id})
		return
	} else if err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while selecting segments"})
		return
	}

	// Вставка нового сегмента
	_, err = tx.Exec(ctx, `INSERT INTO segments (seg_name) VALUES ($1)`, seg.SegName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting segment"})
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
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

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return
	}
	defer tx.Rollback(ctx)

	var segmentID int
	err = tx.QueryRow(ctx, `SELECT seg_id FROM segments WHERE seg_name = $1`, seg.SegName).Scan(&segmentID)

	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while selecting segments"})
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM user_segment WHERE segment_id = $1", segmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting references in user_segment"})
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM segments WHERE seg_name = $1", seg.SegName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting segment"})
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
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

	// Start a transaction
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Add segments to the user
	for _, segName := range us.SegmentsToAdd {

		var segmentID int
		err = tx.QueryRow(ctx, "SELECT seg_id FROM segments WHERE seg_name = $1", segName).Scan(&segmentID)

		log.Println(segmentID)

		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
			continue
		}
		_, err = tx.Exec(ctx, "INSERT INTO user_segment (user_id, segment_id) VALUES ($1, $2)", us.UserID, segmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding segment to user"})
			return
			// should we return ?
		}
	}

	// Remove segments from the user
	for _, segName := range us.SegmentsToRemove {
		var segmentID int
		err = tx.QueryRow(ctx, "SELECT seg_id FROM segments WHERE seg_name = $1", segName).Scan(&segmentID)

		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
			continue
		}

		_, err = tx.Exec(ctx, "DELETE FROM user_segment WHERE user_id = $1 AND segment_id = $2", us.UserID, segmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing segment from user"})
			return
		}
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User segments updated successfully"})
}

// Метод получения активных сегментов пользователя. Принимает на вход id пользователя.

func (s *server) getUserSegments(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("id")
	// var u model.User
	// if err := c.ShouldBindJSON(&u); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
	// 	return
	// }

	rows, err := s.conn.Query(ctx, `
		SELECT seg_id, seg_name
		FROM segments 
		INNER JOIN user_segment ON segments.seg_id = user_segment.segment_id
		WHERE user_segment.user_id = $1 
		`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user segments"})
		return
	}
	defer rows.Close()

	var segments []model.Segment
	for rows.Next() {
		var seg model.Segment
		if err := rows.Scan(&seg.SegID, &seg.SegName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
			return
		}
		segments = append(segments, seg)
	}

	c.JSON(http.StatusOK, segments)
}

func main() {

	dsn := "postgres://postgres:postgres@postgres:5432/postgres" + "?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatalf("can't ping db: %s", err)
	}

	s := server{conn: conn}

	router := gin.Default()

	// router := gin.New()
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	router.GET("/user/:id", s.getUserSegments)
	router.PATCH("/user/:id", s.updateUserSegments)
	router.POST("/segment/create", s.createSegment)
	router.DELETE("/segment/:id", s.deleteSegment)

	log.Fatal(router.Run(":8080"))

	// http.HandleFunc("/create_segment", s.createSegment)
	// http.HandleFunc("/delete_segment", s.deleteSegment)
	// http.HandleFunc("/update_user_segments", s.updateUserSegments)
	// http.HandleFunc("/get_user_segments", s.getUserSegments)

	// log.Fatal(http.ListenAndServe(":8080", nil))
}

// transactions

// func (s *server) deleteSegment(rw http.ResponseWriter, req *http.Request) {
// 	ctx := req.Context()
// 	body, err := io.ReadAll(req.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Println(string(body))

// 	var seg model.SegName

// 	err = json.Unmarshal(body, &seg)
// 	if err != nil {
// 		log.Println("vse horosho")
// 		panic(err)
// 		// make error processing
// 	}

// 	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		http.Error(rw, "Error starting transaction", http.StatusInternalServerError)
// 		return
// 	}
// 	defer tx.Rollback(ctx)
// 	// нужна транзакция
// 	// вынести в папку с базой

// 	var id int
// 	err = tx.QueryRow(ctx, `SELECT seg_id FROM segments WHERE seg_name = $1`, seg.SegName).Scan(&id)

// 	if err == pgx.ErrNoRows {
// 		http.Error(rw, "Error: error while selecting segments: "+err.Error(), http.StatusCreated)
// 		return
// 	}

// 	_, err = tx.Exec(ctx, "DELETE FROM segments WHERE seg_name = $1", seg.SegName)
// 	if err != nil {
// 		http.Error(rw, "Error deleting segment", http.StatusInternalServerError)
// 		return
// 	}

// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		http.Error(rw, "Error committing transaction", http.StatusInternalServerError)
// 		return
// 	}

// 	rw.WriteHeader(http.StatusOK)
// 	rw.Write([]byte("Segment deleted successfully"))
// }

// func (s *server) updateUserSegments(rw http.ResponseWriter, req *http.Request) {
// 	ctx := req.Context()

// 	body, err := io.ReadAll(req.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println(string(body))

// 	var us model.User_segments
// 	err = json.Unmarshal(body, &us)
// 	if err != nil {
// 		log.Println("vse horosho")
// 		panic(err)
// 		// make error processing ??
// 	}
// 	// Start a transaction
// 	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		http.Error(rw, "Error starting transaction", http.StatusInternalServerError)
// 		return
// 	}
// 	defer tx.Rollback(ctx)

// 	// Add segments to the user
// 	for _, segName := range us.SegmentsToAdd {
// 		log.Println(segName)
// 		var segmentID int

// 		err = tx.QueryRow(ctx, "SELECT seg_id FROM segments WHERE seg_name = $1", segName).Scan(&segmentID)

// 		log.Println(segmentID)

// 		if err == pgx.ErrNoRows {
// 			http.Error(rw, "Error: error while selecting segments: "+err.Error(), http.StatusCreated)
// 			continue
// 		}
// 		_, err = tx.Exec(ctx, "INSERT INTO user_segment (user_id, segment_id) VALUES ($1, $2)", us.UserID, segmentID)
// 		if err != nil {
// 			http.Error(rw, "Error adding segment to user: "+err.Error(), http.StatusInternalServerError)
// 		}
// 	}

// 	// Remove segments from the user
// 	for _, segName := range us.SegmentsToRemove {
// 		var segmentID int
// 		err = tx.QueryRow(ctx, "SELECT seg_id FROM segments WHERE seg_name = $1", segName).Scan(&segmentID)

// 		if err == pgx.ErrNoRows {
// 			http.Error(rw, "Error: error while selecting segments: "+err.Error(), http.StatusCreated)
// 			continue
// 		}

// 		_, err = tx.Exec(ctx, "DELETE FROM user_segment WHERE user_id = $1 AND segment_id = $2", us.UserID, segmentID)
// 		if err != nil {
// 			http.Error(rw, "Error removing segment from user: "+err.Error(), http.StatusInternalServerError)
// 		}
// 	}

// 	// Commit the transaction
// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		http.Error(rw, "Error committing transaction", http.StatusInternalServerError)
// 		return
// 	}

// 	rw.WriteHeader(http.StatusOK)
// 	rw.Write([]byte("User segments updated successfully"))
// }
