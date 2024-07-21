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

// ExecuteWithScan used to execute sql query like select
func ExecuteWithScan(cnx *gorm.DB, sql string, params []interface{}) (SelectResult, error) {
	var result SelectResult

	if err := cnx.Raw(sql, params...).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// ExecuteWithExec used to execute sql query like insert, update, delete
func ExecuteWithExec(cnx *gorm.DB, sql string, params []interface{}) error {
	if err := cnx.Exec(sql, params...).Error; err != nil {
		return err
	}

	return nil
}

// GetQuerySQLType used to get the base type (select, insert, update, delete,...) of a sql query
func GetSQLQueryType(sql string) string {
	return strings.ToUpper(strings.SplitN(sql, " ", 2)[0])
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
