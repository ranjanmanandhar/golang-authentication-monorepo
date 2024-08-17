package util

import (
	"errors"
	"os"
)

// GetEnvValue : Get Env Value and Set If Empty

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrForbidden           = errors.New("error forbidden")
	ErrNotFound            = errors.New("item not found")
)

func GetEnvValue(key string, defaultValue string) string {

	value, exists := os.LookupEnv(key)

	if !exists {

		value = defaultValue

	}

	return value

}

func TranslateErrorCode(statuscode int) error {

	switch statuscode {
	case 500:
		return ErrInternalServerError
	case 404:
		return ErrNotFound
	case 403:
		return ErrForbidden
	}
	return nil
}
