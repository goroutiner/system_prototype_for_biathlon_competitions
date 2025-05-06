package entities

import (
	"os"
	"time"
)

// Statistic представляет собой структуру, содержащую статистику участника соревнований по биатлону.
type Statistic struct {
	RequiredStart                 time.Time     // Запланированное время старта
	ActualStart                   time.Time     // Фактическое время старта
	ActualFinish                  time.Time     // Фактическое время финиша
	StartPenaltyLaps              time.Time     // Вспомогательное фактическое время старта прохождения штрафные кругов
	TotalTimeOfPenaltyLaps        time.Duration // Общее время, затраченное на штрафные круги
	TimeOfLapsCompletion          []time.Time   // Временные отметки завершения кругов
	CompetitorID                  string        // Уникальный идентификатор участника
	NumberOfFiringRangeVisited    int           // Количество посещений огневых рубежей
	NumberOfHits                  int           // Количество попаданий в мишени
	NumberOfPenaltyLaps           int           // Количество назначенных штрафных кругов
	NumberOfCompletionPenaltyLaps int           // Количество завершённых штрафных кругов
	NumberOfEndedLaps             int           // Количество завершённых основных кругов
	IsFinished                    bool          // Флаг, указывающий, завершил ли участник гонку
	IsDisqualified                bool          // Флаг, указывающий, был ли участник дисквалифицирован
}

// Files представляет собой структуру, содержащую ссылки на файлы конфигурации и событий.
type Files struct {
	ConfigFile *os.File
	EventsFile *os.File
}

