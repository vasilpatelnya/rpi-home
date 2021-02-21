package jsontool

import (
	"encoding/json"
	"io"
)

func JsonDecode(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)

	return err
}
