package stat

import (
	"api-gateway-sql/utils/uuid"

	"time"
)

// BatchStatistic used to save statistics of all batch requests in memory
type BatchStatistic struct {
	ID             string
	StartTimestamp time.Time
	EndTimestamp   time.Time
	TargetName     string
	SuccessCount   int
	FailureCount   int
	FailureRanges  []string
}

func NewBatchStatistic(targetName string) BatchStatistic {
	return BatchStatistic{
		ID:             uuid.GenerateUID(),
		StartTimestamp: time.Now(),
		TargetName:     targetName,
	}
}
