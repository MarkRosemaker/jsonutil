package jsonutil

import (
	"encoding/json/v2"
	"os"
)

// WriteFile writes a json file by marshalling it.
func WriteFile[T any](name string, data T, opts ...json.Options) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.UnmarshalRead(f, data, opts...); err != nil {
		return err
	}

	return nil
}
