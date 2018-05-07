package check

import "encoding/json"

// ResultPrinter is func declaration that is capable of printing the result
type ResultPrinter func(result ResultInterface) (string, error)

// JsonResultPrinter will encode result into JSON and return it as a string
func JsonResultPrinter(result ResultInterface) (string, error) {
	m, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(m), nil
}
