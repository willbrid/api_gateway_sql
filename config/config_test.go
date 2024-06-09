package config

import (
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
