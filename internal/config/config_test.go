package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/cybre/salesforge-assignment/internal/config"
	"github.com/cybre/salesforge-assignment/internal/database"
)

func TestGetEnv(t *testing.T) {
	// Test case 1: Environment variable exists
	key := "MY_ENV_VAR"
	value := "my_value"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	result := config.GetEnv(key, "fallback")
	if result != value {
		t.Errorf("Expected %s, but got %s", value, result)
	}

	// Test case 2: Environment variable does not exist
	key = "NON_EXISTENT_ENV_VAR"
	fallback := "fallback_value"

	result = config.GetEnv(key, fallback)
	if result != fallback {
		t.Errorf("Expected %s, but got %s", fallback, result)
	}
}

func TestMustGetEnv(t *testing.T) {
	// Test case 1: Environment variable exists
	key := "MY_ENV_VAR"
	value := "my_value"
	os.Setenv(key, value)
	defer os.Unsetenv(key)

	result := config.MustGetEnv(key)
	if result != value {
		t.Errorf("Expected %s, but got %s", value, result)
	}

	// Test case 2: Environment variable does not exist
	key = "NON_EXISTENT_ENV_VAR"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but got none")
		}
	}()

	config.MustGetEnv(key)
}

func TestReadSecret(t *testing.T) {
	// Test case 1: File exists
	filename := "testfile.txt"
	expected := "This is a test secret"
	err := os.WriteFile(filename, []byte(expected), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(filename)

	result := config.ReadSecret(filename)
	os.Remove(filename)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Test case 2: File does not exist
	filename = "nonexistent.txt"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but got none")
		}
	}()

	config.ReadSecret(filename)
}

func TestLoadConfig(t *testing.T) {
	// Test case 1: All environment variables and secret file exist
	os.Setenv("PORT", "4000")
	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_PORT", "5432")
	os.Setenv("DATABASE_NAME", "mydb")
	os.Setenv("DATABASE_USER", "myuser")
	os.Setenv("DATABASE_PASSWORD_FILE", "test_secret.txt")
	os.WriteFile("test_secret.txt", []byte("mysecret"), 0644)

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_HOST")
		os.Unsetenv("DATABASE_PORT")
		os.Unsetenv("DATABASE_NAME")
		os.Unsetenv("DATABASE_USER")
		os.Remove("test_secret.txt")
	}()

	expected := config.Config{
		Port: "4000",
		Database: database.Config{
			Host:     "localhost",
			Port:     "5432",
			Name:     "mydb",
			User:     "myuser",
			Password: "mysecret",
		},
	}

	result := config.LoadConfig()
	os.Remove("test_secret.txt")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, result)
	}
}
