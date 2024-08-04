package stat

import (
	"api-gateway-sql/logging"
	"api-gateway-sql/utils/uuid"

	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(sqlitedb string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s.db", sqlitedb)

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

// BatchStatistic used to save statistics of all batch requests in memory
type BatchStatistic struct {
	ID            string         `json:"id"`
	TargetName    string         `json:"target"`
	SuccessCount  int            `json:"success" gorm:"default:0"`
	FailureCount  int            `json:"failure" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	FailureRanges []FailureRange `json:"failure_ranges"`
}

type FailureRange struct {
	ID        string    `json:"id"`
	StartLine int       `json:"start_line"`
	EndLine   int       `json:"end_line"`
	CreatedAt time.Time `json:"created_at"`
}

func NewBatchStatistic(targetName string) BatchStatistic {
	return BatchStatistic{
		ID:         uuid.GenerateUID(),
		CreatedAt:  time.Now(),
		TargetName: targetName,
	}
}

func AddBatchStatistic(sqlitedb string, bs *BatchStatistic) error {
	cnx, err := Connect(sqlitedb)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	return cnx.Create(bs).Error
}

func UpdateBatchStatistic(sqlitedb string, bs *BatchStatistic, isSuccess bool, failureStartLine int, failureEndLine int) (*BatchStatistic, error) {
	cnx, err := Connect(sqlitedb)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	currentTime := time.Now()
	bs.UpdatedAt = currentTime
	if isSuccess {
		bs.SuccessCount = bs.SuccessCount + 1
	} else {
		bs.FailureCount = bs.FailureCount + 1
		cnx.Model(bs).Association("FailureRanges").Append(&FailureRange{
			ID:        uuid.GenerateUID(),
			StartLine: failureStartLine,
			EndLine:   failureEndLine,
			CreatedAt: currentTime,
		})
	}

	if err := cnx.Save(bs).Error; err != nil {
		return nil, err
	}

	return bs, nil
}

func GetBatchStatistics(sqlitedb string, pageNum int, pageSize int) ([]BatchStatistic, error) {
	var (
		batchStatistics []BatchStatistic
		cnx             *gorm.DB
		err             error
	)

	cnx, err = Connect(sqlitedb)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	err = cnx.Model(&BatchStatistic{}).Preload("FailureRanges").Find(&batchStatistics).Offset(pageNum).Limit(pageSize).Error
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}

	return batchStatistics, err
}
