package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/YugenKawa/dotamatchanalyzer/config"
	"github.com/YugenKawa/dotamatchanalyzer/notifier"
	"github.com/YugenKawa/dotamatchanalyzer/storage"
)

func MatchParse(db *sql.DB, config *config.Config, team config.TeamConfig) {
	log.Println("Starting MatchParse...")

	url := fmt.Sprintf("https://api.opendota.com/api/teams/%d/matches", team.ID)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	var matches []storage.ProMatch
	if err := json.Unmarshal(body, &matches); err != nil {
		fmt.Printf("Error parsing of JSON: %v", err)
		return
	}

	log.Printf("Received %d matches", len(matches))

	if len(matches) == 0 {
		log.Println("No matches found!")
		return
	}

	latestMatches := matches[:min(len(matches), 5)]
	log.Printf("Processing %d latest matches", len(latestMatches))

	for _, match := range latestMatches {
		log.Printf("Cheking match %d", match.ID)

		if storage.IsProcessed(db, match.ID, team.ID) {
			log.Printf("Match %d already processed for team %d, skipping", match.ID, team.ID)
			continue
		}

		message := createMatchMessage(match, team)
		if err := notifier.SendTelegramNotification(config.Telegram.BotToken, config.Telegram.ChatID, message); err != nil {
			log.Printf("Failed to send notification for team %s: %v", team.Name, err)
		}

		if err := storage.SaveProcessedMatch(db, match.ID, team.ID); err != nil {
			log.Printf("Failed to save match %d for team %d: %v", match.ID, team.ID, err)
		} else {
			log.Printf("Successfully saved match %d for team %d", match.ID, team.ID)
		}
	}
}

func createMatchMessage(match storage.ProMatch, team config.TeamConfig) string {
	var playingFor string
	if match.Radiant {
		playingFor = "Radiant"
	} else {
		playingFor = "Dire"
	}

	var result, emoji string
	if (match.Radiant && match.RadiantWin) || (!match.Radiant && !match.RadiantWin) {
		result = fmt.Sprintf("🎉 %s победили!", team.Name)
		emoji = "✅"
	} else {
		result = fmt.Sprintf("💔 %s проиграли", team.Name)
		emoji = "❌"
	}

	startTime := time.Unix(match.StartTime, 0).Format("02.01.2006 15:04")
	duration := match.Duration / 60

	return fmt.Sprintf(
		"%s Новый матч %s!\n\n"+
			"📊 Лига: %s\n"+
			"⏰ Начало: %s\n"+
			"⏱️ Продолжительность: %d мин\n"+
			"📈 Счет: %d - %d\n"+
			"🔗 https://www.opendota.com/matches/%d\n\n"+
			"🏆 %s (%s)",
		emoji,
		team.Name,
		match.LeagueName,
		startTime,
		duration,
		match.RadiantScore,
		match.DireScore,
		match.ID,
		result,
		playingFor,
	)
}
