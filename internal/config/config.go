package config

// Config представляет конфигурацию для системы соревнований по биатлону.
type Config struct {
	Laps        int    `json:"laps"`        // Количество кругов в гонке
	PenaltyLen  int    `json:"penaltyLen"`  // Длина штрафного круга (в метрах)
	LapLen      int    `json:"lapLen"`      // Длина одного круга (в метрах)
	FiringLines int    `json:"firingLines"` // Количество огневых рубежей
	Start       string `json:"start"`       // Время начала гонки
	StartDelta  string `json:"startDelta"`  // Интервал между стартами участников
}
