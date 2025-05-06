package services

import "time"

// formatTimeToDuration преобразует объект времени time.Time в длительность time.Duration.
func formatTimeToDuration(t time.Time) time.Duration {
	refTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	delta := t.Sub(refTime)

	return delta
}

// formatDurationToTime преобразует длительность time.Duration в объект времени time.Time. 
func formatDurationToTime(d time.Duration) time.Time {
	refTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	newTime := refTime.Add(d)

	return newTime
}

