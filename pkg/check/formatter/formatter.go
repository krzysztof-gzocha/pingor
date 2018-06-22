package formatter

import "github.com/krzysztof-gzocha/pingor/pkg/check/result"

// Func is function declaration that is capable of printing the result
type Func func(result result.ResultInterface) (string, error)
