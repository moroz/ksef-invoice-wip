package ksef

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func toJSONPayload(v any) (io.Reader, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("toJSONPayload: %w", err)
	}

	return bytes.NewReader(body), nil
}
