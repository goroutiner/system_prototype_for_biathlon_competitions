package main

import (
	"log"
	"os"
	"system_prototype_for_biathlon_competitions/internal/entities"
	"system_prototype_for_biathlon_competitions/internal/services"
)

func main() {
	configFile, err := os.Open("internal/config/config.json")
	if err != nil {
		log.Printf("failed to open config file: %v\n", err)
		return
	}
	defer configFile.Close()

	eventsFile, err := os.Open("sunny_5_skiers/events")
	if err != nil {
		log.Printf("failed to open events file: %v\n", err)
		return
	}
	defer eventsFile.Close()

	files := &entities.Files{
		ConfigFile: configFile,
		EventsFile: eventsFile,
	}
	service := services.NewParseService(files)

	config, err := service.ParseConfig()
	if err != nil {
		log.Printf("failed to parse config file: %v", err)
		return
	}

	statistics, err := service.ParseEvents(config)
	if err != nil {
		log.Printf("failed to parse events file: %v", err)
		return
	}

	reportService := services.NewReportService(statistics, config)
	if err := reportService.MakeResultingTable(); err != nil {
		log.Printf("failed to make resulting table: %v", err)
		return
	}
}
