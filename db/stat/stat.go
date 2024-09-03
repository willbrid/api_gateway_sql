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
	dsn := fmt.Sprintf("/data/%s.db", sqlitedb)

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

// BatchStatistic used to save statistics of all batch requests in memory
type BatchStatistic struct {
	ID         string    `json:"id"`
	TargetName string    `json:"target"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Blocks     []Block   `json:"blocks" gorm:"foreignKey:BatchStatisticID"`
}

type Block struct {
	ID               string         `json:"id"`
	StartLine        int            `json:"start_line"`
	EndLine          int            `json:"end_line"`
	SuccessCount     int            `json:"success" gorm:"default:0"`
	FailureCount     int            `json:"failure" gorm:"default:0"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	FailureRanges    []FailureRange `json:"failure_ranges" gorm:"foreignKey:BlockID"`
	BatchStatisticID string
}

type FailureRange struct {
	ID        string    `json:"id"`
	StartLine int       `json:"start_line"`
	EndLine   int       `json:"end_line"`
	CreatedAt time.Time `json:"created_at"`
	BlockID   string
}

func NewBatchStatistic(targetName string) *BatchStatistic {
	return &BatchStatistic{
		ID:         uuid.GenerateUID(),
		CreatedAt:  time.Now(),
		TargetName: targetName,
	}
}

func NewBlock(startLine, endLine int) *Block {
	return &Block{
		ID:        uuid.GenerateUID(),
		StartLine: startLine,
		EndLine:   endLine,
		CreatedAt: time.Now(),
	}
}

func NewFailureRange(startLine, endLine int) *FailureRange {
	return &FailureRange{
		ID:        uuid.GenerateUID(),
		StartLine: startLine,
		EndLine:   endLine,
		CreatedAt: time.Now(),
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

func AddNewBlockToBatchStatistic(sqlitedb string, bs *BatchStatistic, startLine int, endLine int) (*Block, error) {
	cnx, err := Connect(sqlitedb)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	bs.UpdatedAt = time.Now()
	block := NewBlock(startLine, endLine)
	err = cnx.Model(bs).Association("Blocks").Append(block)
	if err != nil {
		return nil, err
	}

	if err := cnx.Save(bs).Error; err != nil {
		return nil, err
	}

	return block, nil
}

func UpdateBlock(sqlitedb string, block *Block, isSuccess bool, startLine, endLine int) (*Block, error) {
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
	block.UpdatedAt = currentTime
	if isSuccess {
		block.SuccessCount = block.SuccessCount + 1
	} else {
		block.FailureCount = block.FailureCount + 1
		failureRange := NewFailureRange(startLine, endLine)
		if err := cnx.Model(block).Association("FailureRanges").Append(failureRange); err != nil {
			return nil, err
		}
	}

	if err := cnx.Save(block).Error; err != nil {
		return nil, err
	}

	return block, nil
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

	err = cnx.Model(&BatchStatistic{}).Preload("Blocks").Find(&batchStatistics).Offset(pageNum).Limit(pageSize).Error
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}

	return batchStatistics, err
}
