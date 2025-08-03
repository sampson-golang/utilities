package networking

import (
	"encoding/json"
	"fmt"
)

func UnmarshalResponse(body []byte, data interface{}) error {
	err := json.Unmarshal(body, data)
	if err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}
	return nil
}
