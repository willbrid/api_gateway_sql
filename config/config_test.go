package config_test

import (
	"api-gateway-sql/config"
	"api-gateway-sql/utils/file"

	"fmt"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func TestLoadConfigFileNotFound(t *testing.T) {
	var filename string

	_, err := config.LoadConfig(filename, validate)

	expected := "Config File \"config\" Not Found in \"[]\""

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestLoadConfigFileNotExist(t *testing.T) {
	var filename string = "nonexistentfile.yaml"

	_, err := config.LoadConfig(filename, validate)

	expected := "open nonexistentfile.yaml: no such file or directory"

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestAuthFieldWithAuthEnabled(t *testing.T) {
	configSlices := []string{
		`---
api_gateway_sql:
  auth:
    enabled: true
    username: ""
`,
		`---
api_gateway_sql:
  auth:
    enabled: true
    username: "x"
`,
		`---
api_gateway_sql:
  auth:
    enabled: true
    username: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`,
		`---
api_gateway_sql:
  auth:
    enabled: true
    username: "xxxxx"
    password: ""
`,
		`---
api_gateway_sql:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxx
`,
	}

	expectations := []string{
		"validation failed on field 'Username' for condition 'required_if'",
		"validation failed on field 'Username' for condition 'min'",
		"validation failed on field 'Username' for condition 'max'",
		"validation failed on field 'Password' for condition 'required_if'",
		"validation failed on field 'Password' for condition 'min'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := file.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer os.Remove(filename)

			_, err = config.LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}

func TestSqlitedbField(t *testing.T) {
	configContent := `---
api_gateway_sql:
  sqlitedb: ''
`

	filename, err := file.CreateConfigFileForTesting(configContent)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(filename)

	_, err = config.LoadConfig(filename, validate)

	if err == nil {
		t.Errorf("no error returned")
	}

	if err.Error() == "" {
		t.Errorf("no error message found")
	}
}

func TestDabatasesField(t *testing.T) {
	configSlices := []string{
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "xxxxx"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "1000"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "49152"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: "xxxxx"
    dbname: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    timeout: "10"
`,
	}

	expectations := []string{
		"validation failed on field 'Databases' for condition 'gt'",
		"validation failed on field 'Databases' for condition 'gt'",
		"validation failed on field 'Name' for condition 'required'",
		"validation failed on field 'Name' for condition 'max'",
		"validation failed on field 'Type' for condition 'required'",
		"validation failed on field 'Type' for condition 'oneof'",
		"validation failed on field 'Host' for condition 'required_unless'",
		"validation failed on field 'Host' for condition 'ipv4'",
		"validation failed on field 'Port' for condition 'required_unless'",
		"validation failed on field 'Port' for condition 'min'",
		"validation failed on field 'Port' for condition 'max'",
		"validation failed on field 'Username' for condition 'required_unless'",
		"validation failed on field 'Password' for condition 'required_unless'",
		"validation failed on field 'Dbname' for condition 'required'",
		"validation failed on field 'Timeout' for condition 'required'",
		"1 error(s) decoding:\n\n* error decoding 'api_gateway_sql.databases[0].timeout': time: missing unit in duration \"10\"",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := file.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer os.Remove(filename)

			_, err = config.LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}

func TestTargetsField(t *testing.T) {
	configSlices := []string{
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxx"
    data_source_name: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxx"
    data_source_name: "xxxxx"
    sql: ""
`,
		`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxx"
    data_source_name: "xxxxx"
    sql: "select * from student"
    Multi: true
`,
	}

	expectations := []string{
		"validation failed on field 'Targets' for condition 'gt'",
		"validation failed on field 'Targets' for condition 'gt'",
		"validation failed on field 'Name' for condition 'required'",
		"validation failed on field 'Name' for condition 'max'",
		"validation failed on field 'DataSourceName' for condition 'required'",
		"validation failed on field 'SqlQuery' for condition 'required'",
		"validation failed on field 'BatchSize' for condition 'required_if'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := file.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer os.Remove(filename)

			_, err = config.LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}
