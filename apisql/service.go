package apisql

import (
	"api-gateway-sql/config"
	"api-gateway-sql/db"
	"api-gateway-sql/db/stat"
	"api-gateway-sql/logging"
	"api-gateway-sql/utils/file"

	"fmt"
	"mime/multipart"
	"strings"
	"sync"

	"gorm.io/gorm"
)

// getTargetAndDatabase used to get target and his database
func getTargetAndDatabase(config config.Config, targetName string) (*config.Target, *config.Database, error) {
	target, exist := config.GetTargetByName(targetName)
	if !exist {
		return nil, nil, fmt.Errorf("the specified target name %s does not exist", targetName)
	}

	database, exist := config.GetDatabaseByDataSourceName(target.DataSourceName)
	if !exist {
		return nil, nil, fmt.Errorf("the configured datasource name %s does not exist", target.DataSourceName)
	}

	return &target, &database, nil
}

// executeSingleSQLQuery used to execute a single sql query from target
func executeSingleSQLQuery(target config.Target, database config.Database, postParams map[string]interface{}) (db.SelectResult, error) {
	var (
		cnx    *gorm.DB
		result db.SelectResult = nil
		err    error           = nil
	)

	dbInstance := db.NewDatabase(database.Type)
	cnx, err = dbInstance.Connect(database)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	result, err = dbInstance.ExecuteQuery(target.SqlQuery, postParams)

	return result, err
}

// executeBatchSQLQuery used to execute a batch query from target
func executeBatchSQLQuery(sqlitedb string, target config.Target, database config.Database, postFile multipart.File) error {
	var (
		buffers []file.Buffer
		block   *stat.Block
		bs      *stat.BatchStatistic = stat.NewBatchStatistic(target.Name)
		err     error                = nil
	)

	buffers, err = file.ReadCSVInBuffer(postFile, target.BufferSize)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}

	err = stat.AddBatchStatistic(sqlitedb, bs)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}

	for _, buffer := range buffers {
		block, err = stat.AddNewBlockToBatchStatistic(sqlitedb, bs, buffer.StartLine, buffer.EndLine)
		if err == nil {
			go processBatch(sqlitedb, block, buffer, target, database)
		}
	}

	return nil
}

func processBatch(sqlitedb string, block *stat.Block, buffer file.Buffer, target config.Target, database config.Database) error {
	var (
		cnx *gorm.DB
		err error = nil
	)

	dbInstance := db.NewDatabase(database.Type)
	cnx, err = dbInstance.Connect(database)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	var (
		wg      sync.WaitGroup
		record  map[string]interface{}
		records []map[string]interface{}
	)
	batchSize := target.BatchSize
	batchFields := strings.Split(target.BatchFields, ";")
	numBatches := (len(buffer.Lines) + batchSize - 1) / batchSize

	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(buffer.Lines) {
			end = len(buffer.Lines)
		}

		batch := buffer.Lines[start:end]
		wg.Add(1)

		go func(batch [][]string) {
			defer wg.Done()
			records = make([]map[string]interface{}, len(batch))

			for _, line := range batch {
				record, err = mapBatchFieldToValueLine(batchFields, line)
				if err != nil {
					logging.Log(logging.Error, "%s : %s", err.Error(), line)
					break
				} else {
					records = append(records, record)
				}
			}

			if len(records) > 0 {
				err = dbInstance.ExecuteBatch(target.SqlQuery, records)
				if err != nil {
					logging.Log(logging.Error, err.Error())
					stat.UpdateBlock(sqlitedb, block, false, start, end)
				} else {
					stat.UpdateBlock(sqlitedb, block, true, start, end)
				}
			}
		}(batch)
	}

	wg.Wait()

	return nil
}

// getStats used to get stats in application database
func getStats(sqlitedb string, pageNum int, pageSize int) ([]stat.BatchStatistic, error) {
	var (
		defaultPageNum  int = 1
		defaultPageSize int = 20
	)

	if pageNum > 0 && pageSize > 0 {
		defaultPageNum = pageNum
		defaultPageSize = pageSize
	}

	return stat.GetBatchStatistics(sqlitedb, defaultPageNum, defaultPageSize)
}

// mapBatchFieldToValueLine used to construct a record of a batch sql query
func mapBatchFieldToValueLine(fields []string, values []string) (map[string]interface{}, error) {
	if len(fields) != len(values) {
		return nil, fmt.Errorf("bad mapping fields and file column")
	}

	var result map[string]interface{} = make(map[string]interface{}, len(fields))

	for index, field := range fields {
		result[field] = values[index]
	}

	return result, nil
}

// executeInitSQLQuery used to populate a target database for testing
func executeInitSQLQuery(sqlQuery string, database config.Database) error {
	var (
		cnx *gorm.DB
		err error = nil
	)

	dbInstance := db.NewDatabase(database.Type)
	cnx, err = dbInstance.Connect(database)
	if err != nil {
		logging.Log(logging.Error, err.Error())
		return err
	}
	defer func() {
		dbCnx, _ := cnx.DB()
		dbCnx.Close()
	}()

	queries := strings.Split(sqlQuery, ";")

	return db.ExecuteTransaction(cnx, queries)
}
