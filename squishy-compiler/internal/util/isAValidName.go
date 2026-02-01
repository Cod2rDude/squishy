package util

import (
	"strconv"
	"unicode"

	"github.com/Cod2rDude/squishy/squishy-compiler/internal/errors"
)

// Public Functions
func IsAValidName(str string) (bool, *errors.StackError) {
	if len(str) > 64 {
		return false, errors.New(errors.InvalidNaming, "Length is longer than 64 chars. Length: "+strconv.Itoa(len(str)))
	}

	if len(str) == 0 {
		return false, errors.New(errors.InvalidNaming, "Length is 0.")
	}

	for i, r := range str {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false, errors.New(errors.InvalidNaming, string(r))
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false, errors.New(errors.InvalidNaming, string(r))
			}
		}
	}

	return true, nil
}
