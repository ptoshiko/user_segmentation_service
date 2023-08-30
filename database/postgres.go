package database

import (
	"context"
	"fmt"
	"httpserver/model"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	conn *pgx.Conn
}

func New() (*Database, error) {

	dsn := "postgres://postgres:postgres@postgres:5432/postgres" + "?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		//return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatalf("can't ping db: %s", err)
		//return nil, err
	}

	return &Database{
		conn: conn,
	}, nil
}

func (db *Database) Close(ctx context.Context) {
	db.conn.Close(ctx)
}

func (db *Database) CreateSegment(ctx context.Context, seg model.SegName) error {
	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return err
	}
	defer tx.Rollback(ctx)
	// Проверка, что такого сегмента еще нет
	var id int
	err = tx.QueryRow(ctx, `SELECT seg_id FROM segments WHERE seg_name = $1`, seg.SegName).Scan(&id)

	if err == nil {
		// c.JSON(http.StatusConflict, gin.H{"error": "Segment already exists", "id": id})
		return fmt.Errorf("Segment already exists")
	} else if err != pgx.ErrNoRows {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while selecting segments"})
		return err
	}

	// Вставка нового сегмента
	_, err = tx.Exec(ctx, `INSERT INTO segments (seg_name) VALUES ($1)`, seg.SegName)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting segment"})
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return err
	}
	return nil
}

func (db *Database) DeleteSegment(ctx context.Context, seg model.SegName) error {
	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return err
	}
	defer tx.Rollback(ctx)

	var segmentID int
	err = tx.QueryRow(ctx, `SELECT seg_id FROM segments WHERE seg_name = $1`, seg.SegName).Scan(&segmentID)

	if err == pgx.ErrNoRows {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
		return err
	} else if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while selecting segments"})
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM user_segment WHERE segment_id = $1", segmentID)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting references in user_segment"})
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM segments WHERE seg_name = $1", seg.SegName)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting segment"})
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return err
	}
	return nil
}

// вынести в utils 
func intArrayToString(arr []int) string {
	values := make([]string, len(arr))
	for i, v := range arr {
		values[i] = strconv.Itoa(v)
	}
	return "{" + strings.Join(values, ",") + "}"
}

func (db *Database) UpdateUserSegments(ctx context.Context, us model.UserSegments) error {

	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return err
	}
	defer tx.Rollback(ctx)

	// fmt.Println(us.SegmentsToAdd)

	rows, err := tx.Query(ctx, `
		SELECT seg_id
		FROM segments
		WHERE seg_name = ANY($1)
	`, us.SegmentsToAdd)

	// fmt.Println(err.Error())

	if err != nil {
		return err
	}
	defer rows.Close()

	// var segmentIDs []int
	// for rows.Next() {
	// 	var segmentID int
	// 	if err := rows.Scan(&segmentID); err != nil {
	// 		return err
	// 	}
	// 	segmentIDs = append(segmentIDs, segmentID)
	// }
	// log.Println(segmentIDs)

	// segmentIDsStr := intArrayToString(segmentIDs)

	// _, err = tx.Exec(ctx, `
	// 	INSERT INTO user_segment (user_id, segment_id)
	// 	SELECT $1, unnest($2::integer[])
	// `, us.UserID, segmentIDsStr)
	// if err != nil {
	// 	return err
	// }

	// _, err = tx.Exec(ctx, `
	// 	INSERT INTO user_segment (user_id, segment_id)
	// 	SELECT $1, unnest($2::integer[])
	// `, us.UserID, segmentIDs)

	// // ?? should cast segmentIDs
	// if err != nil {
	// 	return err
	// }

	// // Remove segments from the user
	// for _, segName := range us.SegmentsToRemove {
	// 	var segmentID int
	// 	err = tx.QueryRow(ctx, "SELECT seg_id FROM segments WHERE seg_name = $1", segName).Scan(&segmentID)

	// 	if err == pgx.ErrNoRows {
	// 		// c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
	// 		continue
	// 	}

	// 	_, err = tx.Exec(ctx, "DELETE FROM user_segment WHERE user_id = $1 AND segment_id = $2", us.UserID, segmentID)
	// 	if err != nil {
	// 		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing segment from user"})
	// 		return err
	// 	}
	// }

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return err
	}
	return nil
}

func (db *Database) GetUserSegments(ctx context.Context, userID string) ([]model.Segment, error) {
	rows, err := db.conn.Query(ctx, `
		SELECT seg_id, seg_name
		FROM segments 
		INNER JOIN user_segment ON segments.seg_id = user_segment.segment_id
		WHERE user_segment.user_id = $1 
		`, userID)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user segments"})
		return nil, err
	}
	defer rows.Close()

	var segments []model.Segment
	for rows.Next() {
		var seg model.Segment
		if err := rows.Scan(&seg.SegID, &seg.SegName); err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows"})
			return nil, err
		}
		segments = append(segments, seg)
	}
	return segments, nil
}

// func (db *Database) UpdateUserSegments(ctx context.Context, us model.UserSegments) error {

// 	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	// Add segments to the user

// 	// make the request not in cycle - one request

// 	for _, segName := range us.SegmentsToAdd {

// 		var segmentID int
// 		err = tx.QueryRow(ctx, "SELECT seg_id FROM segments WHERE seg_name = $1", segName).Scan(&segmentID)

// 		log.Println(segmentID)

// 		if err == pgx.ErrNoRows {
// 			// c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
// 			continue
// 		}
// 		_, err = tx.Exec(ctx, "INSERT INTO user_segment (user_id, segment_id) VALUES ($1, $2)", us.UserID, segmentID)
// 		if err != nil {
// 			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding segment to user"})
// 			return err
// 			// should we return ?
// 		}
// 	}

// 	// Remove segments from the user
// 	for _, segName := range us.SegmentsToRemove {
// 		var segmentID int
// 		err = tx.QueryRow(ctx, "SELECT seg_id FROM segments WHERE seg_name = $1", segName).Scan(&segmentID)

// 		if err == pgx.ErrNoRows {
// 			// c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
// 			continue
// 		}

// 		_, err = tx.Exec(ctx, "DELETE FROM user_segment WHERE user_id = $1 AND segment_id = $2", us.UserID, segmentID)
// 		if err != nil {
// 			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing segment from user"})
// 			return err
// 		}
// 	}

// 	// Commit the transaction
// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
// 		return err
// 	}
// 	return nil
// }
