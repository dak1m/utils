package utils

import (
	"encoding/json"
)

type Amount string

var _ json.Marshaler = (*Amount)(nil)

var _ json.Unmarshaler = (*Amount)(nil)

func (a *Amount) UnmarshalJSON(bytes []byte) error {
	var b string
	err := json.Unmarshal(bytes, &b)
	*a = Amount(b)

	return err
}

func (a Amount) MarshalJSON() ([]byte, error) {
	return []byte(FloatStrTrim(string(a))), nil
}
