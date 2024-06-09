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
