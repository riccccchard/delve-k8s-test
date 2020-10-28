package log

import (
	"bytes"
	"fmt"
)

//参数都是偶数个按顺序 key1：value1|key2：value2
func formatw(fields []interface{}) string {
	var buffer bytes.Buffer
	for i := 0; i < len(fields); i++ {
		buffer.WriteString(fmt.Sprint(fields[i]))
		if i%2 == 0 {
			buffer.WriteString(":")
		} else {
			buffer.WriteString("|")
		}
	}
	return buffer.String()
}
