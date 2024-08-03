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
func executeQuery(cnx *gorm.DB, sql string, params []interface{}) (SelectResult, error) {
	var result []map[string]interface{}
	var sqlQueryType string = getSQLQueryType(sql)

	switch sqlQueryType {
	case Select:
		if err := cnx.Raw(sql, params...).Scan(&result).Error; err != nil {
			return nil, err
		}

		return result, nil
	default:
		if err := cnx.Exec(sql, params...).Error; err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// GetQuerySQLType used to get the base type (select, insert, update, delete,...) of a sql query
func getSQLQueryType(sql string) string {
	return strings.ToUpper(strings.SplitN(strings.TrimLeft(sql, " "), " ", 2)[0])
}

// TransformQuery used to parse query from config target
func TransformQuery(query string, params map[string]interface{}) (string, []interface{}) {
	re := regexp.MustCompile(`{{(\w+)}}`)
	matches := re.FindAllStringSubmatch(query, -1)

	values := make([]interface{}, 0, len(matches))
	transformedQuery := re.ReplaceAllStringFunc(query, func(param string) string {
		paramName := param[2 : len(param)-2]
		if value, exists := params[paramName]; exists {
			values = append(values, value)
			return "?"
		}
		return param
	})

	return transformedQuery, values
}
