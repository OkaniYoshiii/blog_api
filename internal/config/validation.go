package config

import (
	"errors"
	"fmt"
	"strconv"
)

type InvalidEnvError struct {
	VarName string
	Cause error
}

func (err *InvalidEnvError) Error() string {
	return fmt.Errorf("invalid environment variable %q : %w", err.VarName, err.Cause).Error()
}

func UndefinedVarError(name string) error {
	return &InvalidEnvError{name, errors.New("variable is not set or empty")}
}

func ValidateDefined(key Key, value string) error {
	if value == "" {
		return UndefinedVarError(key)
	}

	return nil
}

func ParseInt(key Key, value string) (int, error) {
	ttl, err := strconv.Atoi(value)
	if err != nil {
		err := fmt.Errorf("%q cannot be parsed to an integer value: %w", value, err)
		return 0, &InvalidEnvError{"JWT_TTL", err}
	}

	return ttl, nil
}
