package decode

import (
	"io"
	"errors"
	"net/http"
	"encoding/json"
)

var (
	ErrTrailingData = errors.New("request must contain a single JSON object")
)

func JSON[T any](r *http.Request) (T, error) {

	var data T

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&data); err != nil {
		return data, err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
    	return data, ErrTrailingData
	}

	return data, nil
}

