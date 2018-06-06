package printer

import "github.com/krzysztof-gzocha/pingor/pkg/check/result"

// PrinterFunc is func declaration that is capable of printing the result
type PrinterFunc func(result result.ResultInterface) (string, error)
