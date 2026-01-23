package flags

import (
	"bytes"
	"fmt"
)

func MakePretty(raw []byte) string {
	nullIndex := bytes.IndexByte(raw, 0)
	if nullIndex == -1 {
		fmt.Println("Error Couldn't find nullIndex, invalid object format")
		return ""
	}
	content := raw[nullIndex+1:]
	return string(content)
}
