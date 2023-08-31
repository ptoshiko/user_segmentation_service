package database

import (
	"context"
	"fmt"
	"httpserver/model"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
}

func New(dsn string) (*Database, error) {

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Database{
		Conn: conn,
	}, nil
}

func (db *Database) Close(ctx context.Context) {
	db.Conn.Close(ctx)
}

func (db *Database) CreateSegment(ctx context.Context, seg model.SegName) error {
	tx, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx, `SELECT seg_id FROM segments WHERE seg_name = $1`, seg.SegName).Scan(&id)

	if err == nil {
		return fmt.Errorf("Segment already exists")
	} else if err != pgx.ErrNoRows {
		return fmt.Errorf("Segment not found")
	}

	_, err = tx.Exec(ctx, `INSERT INTO segments (seg_name) VALUES ($1)`, seg.SegName)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *Database) DeleteSegment(ctx context.Context, seg model.SegName) error {
	tx, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var segmentID int
	err = tx.QueryRow(ctx, `SELECT seg_id FROM segments WHERE seg_name = $1`, seg.SegName).Scan(&segmentID)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("Segment not found")
	} else if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM user_segment WHERE segment_id = $1", segmentID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM segments WHERE seg_name = $1", seg.SegName)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *Database) UpdateUserSegments(ctx context.Context, us model.UserSegments) error {

	tx, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// Проверка наличия сегментов 
	var segmentIDs []int
	if len(us.SegmentsToAdd) > 0 {
		rows, err := tx.Query(ctx, `
		SELECT seg_id
		FROM segments
		WHERE seg_name = ANY($1)
		`, us.SegmentsToAdd)

		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

		for rows.Next() {
			var segmentID int
			if err := rows.Scan(&segmentID); err != nil {
				return err
			}
			segmentIDs = append(segmentIDs, segmentID)
		}
		rows.Close()

		if len(segmentIDs) != len(us.SegmentsToAdd) {
			return fmt.Errorf("Error in segments_to_add: Segment not found")
		}
	}

	if len(us.SegmentsToRemove) > 0 {
		rows, err := tx.Query(ctx, `
			SELECT segment_id
			FROM user_segment
			WHERE user_id = $1 AND segment_id IN (
				SELECT seg_id
				FROM segments
				WHERE seg_name = ANY($2)
			)
		`, us.UserID, us.SegmentsToRemove)

		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

		var deleteIDs []int
		for rows.Next() {
			var deleteID int
			if err := rows.Scan(&deleteID); err != nil {
				return err
			}
			deleteIDs = append(deleteIDs, deleteID)
		}
		rows.Close()

		if len(deleteIDs) != len(us.SegmentsToRemove) {
			return fmt.Errorf("Error in segments_to_remove: Segment not found")
		}
	}
	// Добавление сегментов
	segmentIDsStr := intArrayToString(segmentIDs)
	_, err = tx.Exec(ctx, `
		INSERT INTO user_segment (user_id, segment_id)
		SELECT $1, unnest($2::integer[])
	`, us.UserID, segmentIDsStr)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// Удаление сегментов 
	query := `
		DELETE FROM user_segment
		WHERE user_id = $1 AND segment_id IN (
			SELECT seg_id FROM segments WHERE seg_name = ANY($2)
		)
	`
	_, err = tx.Exec(ctx, query, us.UserID, us.SegmentsToRemove)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (db *Database) GetUserSegments(ctx context.Context, userID string) ([]model.Segment, error) {
	rows, err := db.Conn.Query(ctx, `
		SELECT seg_id, seg_name
		FROM segments 
		INNER JOIN user_segment ON segments.seg_id = user_segment.segment_id
		WHERE user_segment.user_id = $1 
		`, userID)
	if err != nil {
		return nil, err
	}

	var segments []model.Segment
	for rows.Next() {
		var seg model.Segment
		if err := rows.Scan(&seg.SegID, &seg.SegName); err != nil {
			return nil, err
		}
		segments = append(segments, seg)
	}
	rows.Close()
	return segments, nil
}
