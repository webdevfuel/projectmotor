package test

import (
	"errors"
	"os"
)

func GetTestSession() (string, error) {
	testSession := os.Getenv("TEST_SESSION")
	if testSession == "" {
		return "", errors.New("environment variable TEST_SESSION must be set")
	}
	return testSession, nil
}
