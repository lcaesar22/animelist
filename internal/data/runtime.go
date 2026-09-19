package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// unMarshalJSON splits the runtime request string into two parts and converts the other part to an int64
// MarshalJSON allows custom runtime field(implements the MarshalJSON interface on the runtime)

var errInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r *Runtime) UnmarshalJSON(data []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(data))
	if err != nil {
		return errInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, ".")
	if len(parts) != 2 || parts[1] != "mins" {
		return errInvalidRuntimeFormat
	}

	i, err := strconv.Atoi(parts[0])
	if err != nil {
		return errInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
