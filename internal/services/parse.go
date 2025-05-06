package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"system_prototype_for_biathlon_competitions/internal/config"
	"system_prototype_for_biathlon_competitions/internal/entities"
	"time"
)

// Common format for events:
// [time] eventID competitorID extraParams
//
// Incoming events
// EventID | extraParams | Comments
// 1 	   |             | The competitor registered
// 2       | startTime   | The start time was set by a draw
// 3       |             | The competitor is on the start line
// 4       |             | The competitor has started
// 5       | firingRange | The competitor is on the firing range
// 6 	   | target 	 | The target has been hit
// 7 	   | 			 | The competitor left the firing range
// 8 	   | 			 | The competitor entered the penalty laps
// 9 	   | 			 | The competitor left the penalty laps
// 10 	   | 			 | The competitor ended the main lap
// 11 	   | comment 	 | The competitor can`t continue
//
// Outgoing events
// EventID | extraParams | Comments
// 32 	   | 			 | The competitor is disqualified
// 33 	   | 			 | The competitor has finished

// ParseService представляет сервис для обработки и парсинга файлов.
type ParseService struct {
	files *entities.Files
}

// TimeSet представляет собой структуру, содержащую информацию о времени.
type TimeSet struct {
	ActualTime time.Time
	StartDelta time.Duration
}

func NewParseService(files *entities.Files) *ParseService {
	return &ParseService{files: files}
}

// ParseConfig считывает и парсит конфигурационный файл в формате JSON.
func (s *ParseService) ParseConfig() (*config.Config, error) {
	config := &config.Config{}
	rd := bufio.NewReader(s.files.ConfigFile)
	if err := json.NewDecoder(rd).Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return config, nil
}

// ParseEvents обрабатывает события из файла событий и возвращает статистику участников.
func (s *ParseService) ParseEvents(config *config.Config) (map[string]*entities.Statistic, error) {
	statistics := make(map[string]*entities.Statistic)
	reader := bufio.NewReader(s.files.EventsFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read event: %w", err)
		}

		line = strings.TrimSuffix(line, "\n")
		partition := strings.Split(line, " ")
		if len(partition) < 3 {
			return nil, fmt.Errorf("ivalid incoming events: insufficient number of parameters")
		}

		logTime := partition[0]
		eventID := partition[1]
		competitorID := partition[2]
		actualTime, err := time.Parse("15:04:05.000", strings.Trim(logTime, "[]"))
		if err != nil {
			return nil, fmt.Errorf("ivalid incoming events: failed to parse time: %w", err)
		}

		switch eventID {
		case "1":
			if statistics[competitorID] != nil {
				return nil, fmt.Errorf("ivalid incoming events: competitor has already been registered")
			}
			statistics[competitorID] = &entities.Statistic{}
			statistics[competitorID].CompetitorID = competitorID

			fmt.Printf("%s The competitor(%s) registered\n", logTime, competitorID)
		case "2":
			if len(partition) < 4 {
				return nil, fmt.Errorf("ivalid incoming events: insufficient number of parameters")
			}
			requiredStartStr := partition[3]
			requiredStart, err := time.Parse("15:04:05.000", requiredStartStr)
			if err != nil {
				return nil, fmt.Errorf("ivalid incoming events: failed to parse log time: %w", err)
			}
			statistics[competitorID].RequiredStart = requiredStart

			fmt.Printf("%s The start time for the competitor(%s) was set by a draw to %s\n", logTime, competitorID, requiredStartStr)
		case "3":
			fmt.Printf("%s The competitor(%s) is on the start line\n", logTime, competitorID)
		case "4":
			statistics[competitorID].ActualStart = actualTime
			deltaTime, err := time.Parse("03:04:05", config.StartDelta)
			if err != nil {
				return nil, fmt.Errorf("ivalid incoming events: failed to parse delta time in config: %w", err)
			}

			startDelta := formatTimeToDuration(deltaTime)
			timeSet := &TimeSet{
				ActualTime: actualTime,
				StartDelta: startDelta,
			}
			DisqualifiedCheck(statistics, timeSet)

			fmt.Printf("%s The competitor(%s) has started\n", logTime, competitorID)
		case "5":
			if len(partition) < 4 {
				return nil, fmt.Errorf("ivalid incoming events: insufficient number of parameters")
			}
			firingRange := partition[3]
			statistics[competitorID].NumberOfFiringRangeVisited++
			statistics[competitorID].NumberOfPenaltyLaps = 5

			fmt.Printf("%s The competitor(%s) is on the firing range(%s)\n", logTime, competitorID, firingRange)
		case "6":
			if len(partition) < 4 {
				return nil, fmt.Errorf("ivalid incoming events: insufficient number of parameters")
			}
			target := partition[3]
			statistics[competitorID].NumberOfHits++
			statistics[competitorID].NumberOfPenaltyLaps--

			fmt.Printf("%s The target(%s) has been hit by competitor(%s)\n", logTime, target, competitorID)
		case "7":
			fmt.Printf("%s The competitor(%s) left the firing range\n", logTime, competitorID)
		case "8":
			statistics[competitorID].StartPenaltyLaps = actualTime

			fmt.Printf("%s The competitor(%s) entered the penalty laps\n", logTime, competitorID)
		case "9":
			startPenaltyLaps := statistics[competitorID].StartPenaltyLaps
			statistics[competitorID].NumberOfCompletionPenaltyLaps += statistics[competitorID].NumberOfPenaltyLaps
			statistics[competitorID].TotalTimeOfPenaltyLaps += actualTime.Sub(startPenaltyLaps)

			fmt.Printf("%s The competitor(%s) left the penalty laps\n", logTime, competitorID)
		case "10":
			statistics[competitorID].NumberOfEndedLaps++
			statistics[competitorID].TimeOfLapsCompletion = append(statistics[competitorID].TimeOfLapsCompletion, actualTime)
			if statistics[competitorID].NumberOfEndedLaps != config.Laps {
				fmt.Printf("%s The competitor(%s) ended the main lap\n", logTime, competitorID)
				continue
			}
			statistics[competitorID].IsFinished = true
			statistics[competitorID].ActualFinish = actualTime

			fmt.Printf("%s The competitor(%s) has finished\n", logTime, competitorID)
		case "11":
			comment := strings.Join(partition[3:], " ")

			fmt.Printf("%s The competitor(%s) can`t continue: %s\n", logTime, competitorID, comment)
		}
	}

	return statistics, nil
}

// DisqualifiedCheck проверяет, был ли участник дисквалифицирован на основе
// предоставленной статистики и текущего времени.
func DisqualifiedCheck(statistics map[string]*entities.Statistic, timeSet *TimeSet) {
	for competitorID, statistic := range statistics {
		requiredStart := statistic.RequiredStart

		if requiredStart.IsZero() || statistic.IsDisqualified || !statistic.ActualStart.IsZero() {
			continue
		}
		if timeSet.ActualTime.After(requiredStart.Add(timeSet.StartDelta)) {
			actualTimeStr := timeSet.ActualTime.Format("15:04:05.000")
			fmt.Printf("%s The competitor(%s) is disqualified\n", actualTimeStr, competitorID)
			statistics[competitorID].IsDisqualified = true
		}
	}
}
