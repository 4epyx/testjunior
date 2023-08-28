package util

import (
	"fmt"
	"os"
)

// ParseEnv parses the environment with names given in the varNames. It returns map representation of environment variables,
// where key is name of the variable and value is the value of the variable and slice of errors for every environment variable, that not found.
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
