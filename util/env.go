package util

import (
	"fmt"
	"os"
)

func ParseEnv(varNames ...string) (map[string]string, []error) {
	res := make(map[string]string)
	errors := make([]error, 0)

	for _, varName := range varNames {
		val, ok := os.LookupEnv(varName)
		if !ok {
			errors = append(errors, fmt.Errorf("couldn't find environment variable with name %s", varName))
		}

		res[varName] = val
	}

	return res, errors
}
