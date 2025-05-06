package services

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
	"system_prototype_for_biathlon_competitions/internal/config"
	"system_prototype_for_biathlon_competitions/internal/entities"
	"time"
)

// Format:
// - Total time includes the difference between scheduled and actual start time or NotStarted/NotFinished marks
// - Time taken to complete each lap
// - Average speed for each lap [m/s]
// - Time taken to complete penalty laps
// - Average speed over penalty laps [m/s]
// - Number of hits/number of shots
//
// Example:
// [NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/5

// ReportService предоставляет сервис для работы с отчетами.
type ReportService struct {
	Statistics map[string]*entities.Statistic
	Config     *config.Config
}

func NewReportService(statistics map[string]*entities.Statistic, config *config.Config) *ReportService {
	return &ReportService{
		Statistics: statistics,
		Config:     config,
	}
}

// MakeResultingTable создает итоговую таблицу результатов соревнований и записывает её в файл 'report'.
func (s *ReportService) MakeResultingTable() error {
	sortedStatistics := s.SortStatistics()
	reportFile, err := os.Create("report")
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer reportFile.Close()

	writer := bufio.NewWriter(reportFile)
	defer writer.Flush()

	for _, statistic := range sortedStatistics {
		totalTime := GetTotalTime(statistic)
		timeAndAvgSpeedForLaps := GetTimeAndAvgSpeedForLaps(statistic, s.Config)
		timeAndAvgSpeedForPenaltyLaps := GetTimeAndAvgSpeedForPenaltyLaps(statistic, s.Config)
		hitStatistics := GetHitStatistics(statistic)
		line := fmt.Sprintf("[%s] %s [%s] %s %s\n", totalTime, statistic.CompetitorID, timeAndAvgSpeedForLaps, timeAndAvgSpeedForPenaltyLaps, hitStatistics)
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write line in report file: %w", err)
		}
	}

	return nil
}

// SortStatistics сортирует статистику участников соревнований в три категории:
// завершившие, незавершившие и дисквалифицированные. 
func (s *ReportService) SortStatistics() []*entities.Statistic {
	finishedList := make([]*entities.Statistic, 0, len(s.Statistics))
	notFinishedList := make([]*entities.Statistic, 0, len(s.Statistics))
	disqualifiedList := make([]*entities.Statistic, 0, len(s.Statistics))
	for _, statistic := range s.Statistics {
		if statistic.IsFinished {
			finishedList = append(finishedList, statistic)
		} else if statistic.IsDisqualified {
			disqualifiedList = append(disqualifiedList, statistic)
		} else {
			notFinishedList = append(notFinishedList, statistic)
		}
	}

	slices.SortFunc(finishedList, func(a, b *entities.Statistic) int {
		if a.ActualFinish.Sub(a.RequiredStart) > b.ActualFinish.Sub(b.RequiredStart) {
			return 1
		} else if a.ActualFinish.Sub(a.RequiredStart) < b.ActualFinish.Sub(b.RequiredStart) {
			return -1
		}
		return 0
	})
	slices.SortFunc(notFinishedList, func(a, b *entities.Statistic) int {
		return cmp.Compare(a.CompetitorID, b.CompetitorID)
	})
	slices.SortFunc(disqualifiedList, func(a, b *entities.Statistic) int {
		return cmp.Compare(a.CompetitorID, b.CompetitorID)
	})

	sortedStatistics := append(append(finishedList, notFinishedList...), disqualifiedList...)

	return sortedStatistics
}

// GetTotalTime вычисляет общее время прхождения эстафеты на основе статистики участника.
func GetTotalTime(statistic *entities.Statistic) string {
	if statistic.IsDisqualified {
		return "NotStarted"
	}
	if !statistic.IsFinished {
		return "NotFinished"
	}

	totalInterval := statistic.ActualFinish.Sub(statistic.RequiredStart)
	totalTime := formatDurationToTime(totalInterval)
	totalTimeStr := totalTime.Format("15:04:05.000")

	return totalTimeStr
}

// GetTimeAndAvgSpeedForLaps вычисляет время прохождения и среднюю скорость для каждого круга.
func GetTimeAndAvgSpeedForLaps(statistic *entities.Statistic, config *config.Config) string {
	distance := float64(config.LapLen)
	pairs := make([]string, config.Laps)
	for i := 0; i < config.Laps; i++ {
		var (
			takenInterval time.Duration
			takenTime     time.Time
		)
		if i < len(statistic.TimeOfLapsCompletion) {
			if i == 0 {
				takenInterval = statistic.TimeOfLapsCompletion[i].Sub(statistic.ActualStart)
				takenTime = formatDurationToTime(takenInterval)
			} else {
				takenInterval = statistic.TimeOfLapsCompletion[i].Sub(statistic.TimeOfLapsCompletion[i-1])
				takenTime = formatDurationToTime(takenInterval)
			}
		}

		var pair string
		sec := takenInterval.Seconds()
		takenTimeStr := takenTime.Format("15:04:05.000")
		if sec == 0 {
			pairs[i] = "{,}"
			continue
		}
		avgSpeed := distance / sec
		pair = fmt.Sprintf("{%s, %.3f}", takenTimeStr, avgSpeed)
		pairs[i] = pair
	}
	mergedPairs := strings.Join(pairs, ", ")

	return mergedPairs
}

// GetTimeAndAvgSpeedForPenaltyLaps вычисляет общее время и среднюю скорость для штрафных кругов.
func GetTimeAndAvgSpeedForPenaltyLaps(statistic *entities.Statistic, config *config.Config) string {
	var pair string
	distance := float64(config.PenaltyLen*statistic.NumberOfCompletionPenaltyLaps) * 1000
	totalInterval := statistic.TotalTimeOfPenaltyLaps
	msec := float64(totalInterval.Milliseconds())

	if msec == 0 {
		return "{,}"
	}
	avgSpeed := distance / msec
	totalTime := formatDurationToTime(totalInterval)
	totalTimeStr := totalTime.Format("15:04:05.000")
	pair = fmt.Sprintf("{%s, %.3f}", totalTimeStr, avgSpeed)

	return pair
}

// GetHitStatistics возвращает строковое представление статистики попаданий
// в формате "количество попаданий/общее количество выстрелов".
func GetHitStatistics(statistic *entities.Statistic) string {
	numberOfHits := statistic.NumberOfHits
	numberOfShots := 5 * statistic.NumberOfFiringRangeVisited
	stat := fmt.Sprintf("%d/%d", numberOfHits, numberOfShots)

	return stat
}
