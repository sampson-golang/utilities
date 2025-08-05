package output

import (
	"encoding/json"
)

func PrettifyBytes(i interface{}, indent ...string) []byte {
	indentStr := "  "
	if len(indent) > 0 {
		indentStr = indent[0]
	}
	bytes, _ := json.MarshalIndent(i, "", indentStr)
	return bytes
}

func Prettify(i interface{}, indent ...string) string {
	return string(PrettifyBytes(i, indent...))
}
