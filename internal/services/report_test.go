package services_test

import (
	"system_prototype_for_biathlon_competitions/internal/config"
	"system_prototype_for_biathlon_competitions/internal/entities"
	"system_prototype_for_biathlon_competitions/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestSortStatistics тестирует метод SortStatistics.
func TestSortStatistics(t *testing.T) {
	tests := []struct {
		name       string
		statistics map[string]*entities.Statistic
		expected   []*entities.Statistic
	}{
		{
			name: "All competitors finished",
			statistics: map[string]*entities.Statistic{
				"1": {
					CompetitorID:  "1",
					IsFinished:    true,
					RequiredStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
					ActualFinish:  time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
				},
				"2": {
					CompetitorID:  "2",
					IsFinished:    true,
					RequiredStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
					ActualFinish:  time.Date(2023, 10, 1, 10, 25, 0, 0, time.UTC),
				},
			},
			expected: []*entities.Statistic{
				{
					CompetitorID:  "2",
					IsFinished:    true,
					RequiredStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
					ActualFinish:  time.Date(2023, 10, 1, 10, 25, 0, 0, time.UTC),
				},
				{
					CompetitorID:  "1",
					IsFinished:    true,
					RequiredStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
					ActualFinish:  time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Mixed competitors",
			statistics: map[string]*entities.Statistic{
				"1": {
					CompetitorID:  "1",
					IsFinished:    true,
					RequiredStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
					ActualFinish:  time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
				},
				"2": {
					CompetitorID:   "2",
					IsDisqualified: true,
				},
				"3": {
					CompetitorID: "3",
					IsFinished:   false,
				},
			},
			expected: []*entities.Statistic{
				{
					CompetitorID:  "1",
					IsFinished:    true,
					RequiredStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
					ActualFinish:  time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
				},
				{
					CompetitorID: "3",
					IsFinished:   false,
				},
				{
					CompetitorID:   "2",
					IsDisqualified: true,
				},
			},
		},
		{
			name: "All competitors disqualified",
			statistics: map[string]*entities.Statistic{
				"2": {
					CompetitorID:   "2",
					IsDisqualified: true,
				},
				"1": {
					CompetitorID:   "1",
					IsDisqualified: true,
				},
			},
			expected: []*entities.Statistic{
				{
					CompetitorID:   "1",
					IsDisqualified: true,
				},
				{
					CompetitorID:   "2",
					IsDisqualified: true,
				},
			},
		},
	}

	for _, tt := range tests {
		service := &services.ReportService{
			Statistics: tt.statistics,
		}
		sortedStatistics := service.SortStatistics()
		require.Equal(t, len(tt.expected), len(sortedStatistics))
		for i, expectedStat := range tt.expected {
			require.Equal(t, expectedStat.CompetitorID, sortedStatistics[i].CompetitorID)
		}
	}
}

// TestGetTotalTime тестирует функцию GetTotalTime.
func TestGetTotalTime(t *testing.T) {
	tests := []struct {
		name      string
		statistic *entities.Statistic
		expected  string
	}{
		{
			name: "Disqualified competitor",
			statistic: &entities.Statistic{
				IsDisqualified: true,
			},
			expected: "NotStarted",
		},
		{
			name: "Not finished competitor",
			statistic: &entities.Statistic{
				IsFinished: false,
			},
			expected: "NotFinished",
		},
		{
			name: "Finished competitor",
			statistic: &entities.Statistic{
				IsFinished:    true,
				RequiredStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				ActualFinish:  time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
			},
			expected: "00:30:00.000",
		},
	}

	for _, tt := range tests {
		result := services.GetTotalTime(tt.statistic)
		require.Equal(t, tt.expected, result)
	}
}

// TestGetTimeAndAvgSpeedForLaps тестирует функцию GetTimeAndAvgSpeedForLaps.
func TestGetTimeAndAvgSpeedForLaps(t *testing.T) {
	tests := []struct {
		name      string
		statistic *entities.Statistic
		config    *config.Config
		expected  string
	}{
		{
			name: "No laps completed",
			statistic: &entities.Statistic{
				ActualStart:          time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				TimeOfLapsCompletion: []time.Time{},
			},
			config: &config.Config{
				LapLen: 1000,
				Laps:   3,
			},
			expected: "{,}, {,}, {,}",
		},
		{
			name: "One lap completed",
			statistic: &entities.Statistic{
				ActualStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				TimeOfLapsCompletion: []time.Time{
					time.Date(2023, 10, 1, 10, 5, 0, 0, time.UTC),
				},
			},
			config: &config.Config{
				LapLen: 1000,
				Laps:   3,
			},
			expected: "{00:05:00.000, 3.333}, {,}, {,}",
		},
		{
			name: "All laps completed",
			statistic: &entities.Statistic{
				ActualStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				TimeOfLapsCompletion: []time.Time{
					time.Date(2023, 10, 1, 10, 5, 0, 0, time.UTC),
					time.Date(2023, 10, 1, 10, 10, 0, 0, time.UTC),
					time.Date(2023, 10, 1, 10, 15, 0, 0, time.UTC),
				},
			},
			config: &config.Config{
				LapLen: 1000,
				Laps:   3,
			},
			expected: "{00:05:00.000, 3.333}, {00:05:00.000, 3.333}, {00:05:00.000, 3.333}",
		},
		{
			name: "Partial laps completed",
			statistic: &entities.Statistic{
				ActualStart: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				TimeOfLapsCompletion: []time.Time{
					time.Date(2023, 10, 1, 10, 5, 0, 0, time.UTC),
					time.Date(2023, 10, 1, 10, 12, 0, 0, time.UTC),
				},
			},
			config: &config.Config{
				LapLen: 1000,
				Laps:   3,
			},
			expected: "{00:05:00.000, 3.333}, {00:07:00.000, 2.381}, {,}",
		},
	}

	for _, tt := range tests {
		result := services.GetTimeAndAvgSpeedForLaps(tt.statistic, tt.config)
		require.Equal(t, tt.expected, result)
	}
}
// TestGetTimeAndAvgSpeedForPenaltyLaps тестирует функцию GetTimeAndAvgSpeedForPenaltyLaps.
func TestGetTimeAndAvgSpeedForPenaltyLaps(t *testing.T) {
	tests := []struct {
		name      string
		statistic *entities.Statistic
		config    *config.Config
		expected  string
	}{
		{
			name: "No penalty laps completed",
			statistic: &entities.Statistic{
				TotalTimeOfPenaltyLaps:        0,
				NumberOfCompletionPenaltyLaps: 0,
			},
			config: &config.Config{
				PenaltyLen: 150,
			},
			expected: "{,}",
		},
		{
			name: "One penalty lap completed",
			statistic: &entities.Statistic{
				TotalTimeOfPenaltyLaps:        2 * time.Minute,
				NumberOfCompletionPenaltyLaps: 1,
			},
			config: &config.Config{
				PenaltyLen: 150,
			},
			expected: "{00:02:00.000, 1.250}",
		},
		{
			name: "Multiple penalty laps completed",
			statistic: &entities.Statistic{
				TotalTimeOfPenaltyLaps:        5 * time.Minute,
				NumberOfCompletionPenaltyLaps: 3,
			},
			config: &config.Config{
				PenaltyLen: 150,
			},
			expected: "{00:05:00.000, 1.500}",
		},
	}

	for _, tt := range tests {
        result := services.GetTimeAndAvgSpeedForPenaltyLaps(tt.statistic, tt.config)
        require.Equal(t, tt.expected, result)
	}
}
