package config

import (
	"db-api-sql/utils"

	"fmt"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func TestLoadConfigFileNotFound(t *testing.T) {
	var filename string

	_, err := LoadConfig(filename, validate)

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

	_, err := LoadConfig(filename, validate)

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
db_api_sql:
  auth:
    enabled: true
    username: ""
`,
		`---
db_api_sql:
  auth:
    enabled: true
    username: "x"
`,
		`---
db_api_sql:
  auth:
    enabled: true
    username: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`,
		`---
db_api_sql:
  auth:
    enabled: true
    username: "xxxxx"
    password: ""
`,
		`---
db_api_sql:
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
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

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

func TestTimeoutField(t *testing.T) {
	configContent := `---
db_api_sql:
  timeout: ''
`

	filename, err := utils.CreateConfigFileForTesting(configContent)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.Remove(filename)

	_, err = LoadConfig(filename, validate)

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
db_api_sql:
  timeout: "10s"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
`,
		`---
db_api_sql:
  timeout: "10s"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
`,
		`---
db_api_sql:
  timeout: "10s"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: ""
`,
		`---
db_api_sql:
  timeout: "10s"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`,
		`---
db_api_sql:
  timeout: "10s"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: ""
`,
		`---
db_api_sql:
  timeout: "10s"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "xxxxx"
`,
		`---
db_api_sql:
  timeout: "10s"
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
db_api_sql:
  timeout: "10s"
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
db_api_sql:
  timeout: "10s"
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
db_api_sql:
  timeout: "10s"
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
db_api_sql:
  timeout: "10s"
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
db_api_sql:
  timeout: "10s"
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
db_api_sql:
  timeout: "10s"
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
db_api_sql:
  timeout: "10s"
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
	}

	expectations := []string{
		"validation failed on field 'Databases' for condition 'gt'",
		"validation failed on field 'Databases' for condition 'gt'",
		"validation failed on field 'Name' for condition 'required'",
		"validation failed on field 'Name' for condition 'max'",
		"validation failed on field 'Type' for condition 'required'",
		"validation failed on field 'Type' for condition 'oneof'",
		"validation failed on field 'Host' for condition 'required'",
		"validation failed on field 'Host' for condition 'ipv4'",
		"validation failed on field 'Port' for condition 'required'",
		"validation failed on field 'Port' for condition 'min'",
		"validation failed on field 'Port' for condition 'max'",
		"validation failed on field 'Username' for condition 'required'",
		"validation failed on field 'Password' for condition 'required'",
		"validation failed on field 'Dbname' for condition 'required'",
		"",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

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
