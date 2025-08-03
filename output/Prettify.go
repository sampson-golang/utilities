package output

import (
	"encoding/json"
)

func Prettify(i interface{}, indent ...string) string {
	indentStr := "  "
	if len(indent) > 0 {
		indentStr = indent[0]
	}
	stringified, _ := json.MarshalIndent(i, "", indentStr)
	return string(stringified)
}
