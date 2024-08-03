package db

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type SelectResult []map[string]interface{}

const (
	Select string = "SELECT"
)

// executeQuery used to execute a simple non batch query
func executeQuery(cnx *gorm.DB, sqlQuery string, params map[string]interface{}) (SelectResult, error) {
	var result []map[string]interface{}
	var sqlQueryType string = getSQLQueryType(sqlQuery)
	parsedSql, parsedParams := transformQuery(sqlQuery, params)

	switch sqlQueryType {
	case Select:
		if err := cnx.Raw(parsedSql, parsedParams...).Scan(&result).Error; err != nil {
			return nil, err
		}

		return result, nil
	default:
		if err := cnx.Exec(parsedSql, parsedParams...).Error; err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// executeBatch used to batch sql query
func executeBatch(cnx *gorm.DB, sqlQuery string, params []map[string]interface{}) error {
	return cnx.Transaction(func(tx *gorm.DB) error {
		for _, param := range params {
			parsedSql, parsedParams := transformQuery(sqlQuery, param)
			if err := tx.Exec(parsedSql, parsedParams...).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// getQuerySQLType used to get the base type (select, insert, update, delete,...) of a sql query
func getSQLQueryType(sqlQuery string) string {
	return strings.ToUpper(strings.SplitN(strings.TrimLeft(sqlQuery, " "), " ", 2)[0])
}

// transformQuery used to parse query from config target
func transformQuery(sqlQuery string, params map[string]interface{}) (string, []interface{}) {
	re := regexp.MustCompile(`{{(\w+)}}`)
	matches := re.FindAllStringSubmatch(sqlQuery, -1)

	values := make([]interface{}, 0, len(matches))
	transformedQuery := re.ReplaceAllStringFunc(sqlQuery, func(param string) string {
		paramName := param[2 : len(param)-2]
		if value, exists := params[paramName]; exists {
			values = append(values, value)
			return "?"
		}
		return param
	})

	return transformedQuery, values
}
