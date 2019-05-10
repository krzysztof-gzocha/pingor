package json

import (
	"encoding/json"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// Formatter will encode result into JSON and return it as a string
func Formatter(result result.Result) (string, error) {
	m, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(m), nil
}
