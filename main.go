package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YugenKawa/dotamatchanalyzer/api"
	"github.com/YugenKawa/dotamatchanalyzer/config"
	"github.com/YugenKawa/dotamatchanalyzer/storage"
)

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db := storage.InitDB(&config.Database)
	defer db.Close()

	log.Println("✅ PostgreSQL подключен успешно!")

	for _, team := range config.Teams {
		log.Printf("Starting monitoring for team: %s (ID: %d)", team.Name, team.ID)
		go monitorTeam(db, config, team)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("🚀 Сервис запущен! Ожидаем сигнал завершения (Ctrl+C)...")
	<-sigChan

	log.Println("🛑Завершаем работу...")
}

func monitorTeam(db *sql.DB, config *config.Config, team config.TeamConfig) {
	interval, err := time.ParseDuration(config.Settings.CheckInterval)
	if err != nil {
		log.Printf("❌Invalid interval for team %s: %v", team.Name, err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	api.MatchParse(db, config, team)

	log.Printf("🔁 Monitoring team %s with interval %s", team.Name, interval)

	for {
		select {
		case <-ticker.C:
			api.MatchParse(db, config, team)
		}
	}
}
