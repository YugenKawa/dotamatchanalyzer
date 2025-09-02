package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/YugenKawa/dotamatchanalyzer/config"
	_ "github.com/lib/pq"
)

type ProMatch struct {
	ID           int64  `json:"match_id"`
	Radiant      bool   `json:"radiant"`
	RadiantWin   bool   `json:"radiant_win"`
	RadiantScore int    `json:"radiant_score"`
	DireScore    int    `json:"dire_score"`
	LeagueID     int    `json:"leagueid"`
	LeagueName   string `json:"league_name"`
	Duration     int    `json:"duration"`
	StartTime    int64  `json:"start_time"`
}

type ProTeam struct {
	ID            int    `json:"account_id"`
	Name          string `json:"name"`
	GamePlayed    int    `json:"games_played"`
	Wins          int    `json:"wins"`
	CurrentMember bool   `json:"is_current_team_member"`
}

func InitDB(dbConfig *config.DatabaseConfig) *sql.DB {
	dsn := dbConfig.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	log.Println("✅Successfully connected to PostgreSQL!")

	createTable(db)

	return db
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS processed_matches (
	match_id BIGINT,
	team_id INTEGER,
	processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (match_id, team_id)
	);

	CREATE INDEX IF NOT EXISTS idx_processed_matches_team_id
	ON processed_matches(team_id);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	var columnExists bool
	checkQuery := `
	SELECT EXISTS (
	SELECT 1
	FROM information_schema.columns
	WHERE table_name = 'processed_matches'
	AND column_name = 'team_id'
	)`

	err = db.QueryRow(checkQuery).Scan(&columnExists)
	if err != nil {
		log.Printf("Error cheking for team_id column: %v", err)
		return
	}

	if !columnExists {
		log.Println("Adding team_id column to existing table...")
		alterQuery := `ALTER TABLE processed_matches ADD COLUMN team_id INTEGER`
		_, err := db.Exec(alterQuery)
		if err != nil {
			log.Fatal("Failed to add team_id column:", err)
		}
		log.Println("✅team_id column added successfully!")
	}
}
