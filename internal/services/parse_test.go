package services_test

import (
	"testing"
	"time"

	"system_prototype_for_biathlon_competitions/internal/entities"
	"system_prototype_for_biathlon_competitions/internal/services"

	"github.com/stretchr/testify/require"
)

// TestDisqualifiedCheck тестирует функцию DisqualifiedCheck.
func TestDisqualifiedCheck(t *testing.T) {
	t.Run("competitor is disqualified", func(t *testing.T) {
		statistics := map[string]*entities.Statistic{
			"1": {
				CompetitorID:   "1",
				RequiredStart:  time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				IsDisqualified: false,
				ActualStart:    time.Time{},
			},
		}
		timeSet := &services.TimeSet{
			ActualTime: time.Date(2023, 10, 1, 10, 5, 0, 0, time.UTC),
			StartDelta: 3 * time.Minute,
		}

		services.DisqualifiedCheck(statistics, timeSet)

		require.True(t, statistics["1"].IsDisqualified)
	})

	t.Run("start on time", func(t *testing.T) {
		statistics := map[string]*entities.Statistic{
			"1": {
				CompetitorID:   "1",
				RequiredStart:  time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				IsDisqualified: false,
				ActualStart:    time.Date(2023, 10, 1, 10, 2, 0, 0, time.UTC),
			},
		}
		timeSet := &services.TimeSet{
			ActualTime: time.Date(2023, 10, 1, 10, 2, 0, 0, time.UTC),
			StartDelta: 3 * time.Minute,
		}

		services.DisqualifiedCheck(statistics, timeSet)

		require.False(t, statistics["1"].IsDisqualified)
	})

	t.Run("competitor is already disqualified", func(t *testing.T) {
		statistics := map[string]*entities.Statistic{
			"1": {
				CompetitorID:   "1",
				RequiredStart:  time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				IsDisqualified: true,
				ActualStart:    time.Time{},
			},
		}
		timeSet := &services.TimeSet{
			ActualTime: time.Date(2023, 10, 1, 10, 5, 0, 0, time.UTC),
			StartDelta: 3 * time.Minute,
		}

		services.DisqualifiedCheck(statistics, timeSet)

		require.True(t, statistics["1"].IsDisqualified)
	})

	t.Run("actual start time is already set", func(t *testing.T) {
		statistics := map[string]*entities.Statistic{
			"1": {
				CompetitorID:   "1",
				RequiredStart:  time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
				IsDisqualified: false,
				ActualStart:    time.Date(2023, 10, 1, 10, 1, 0, 0, time.UTC),
			},
		}
		timeSet := &services.TimeSet{
			ActualTime: time.Date(2023, 10, 1, 10, 5, 0, 0, time.UTC),
			StartDelta: 3 * time.Minute,
		}

		services.DisqualifiedCheck(statistics, timeSet)

		require.False(t, statistics["1"].IsDisqualified)
	})

	t.Run("required start time is zero", func(t *testing.T) {
		statistics := map[string]*entities.Statistic{
			"1": {
				CompetitorID:   "1",
				RequiredStart:  time.Time{},
				IsDisqualified: false,
				ActualStart:    time.Time{},
			},
		}
		timeSet := &services.TimeSet{
			ActualTime: time.Date(2023, 10, 1, 10, 5, 0, 0, time.UTC),
			StartDelta: 3 * time.Minute,
		}

		services.DisqualifiedCheck(statistics, timeSet)

		require.False(t, statistics["1"].IsDisqualified)
	})
}
