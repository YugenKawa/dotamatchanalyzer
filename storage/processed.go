package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func IsProcessed(db *sql.DB, matchID int64, teamID int) bool {
	log.Printf("Cheking if match %d is processed", matchID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM processed_matches WHERE match_id = $1 AND team_id = $2)`
	err := db.QueryRowContext(ctx, query, matchID, teamID).Scan(&exists)
	if err != nil {
		log.Printf("Error cheking match %d for team %d: %v", matchID, teamID, err)
		return false
	}

	return exists
}

func SaveProcessedMatch(db *sql.DB, matchID int64, teamID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO processed_matches (match_id, team_id) VALUES ($1, $2)`
	_, err := db.ExecContext(ctx, query, matchID, teamID)
	if err != nil {
		return fmt.Errorf("failed to save match %d for team %d: %v", matchID, teamID, err)
	}
	return nil
}
