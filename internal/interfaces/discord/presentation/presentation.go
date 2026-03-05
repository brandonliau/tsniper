package presentation

import (
	"fmt"
)

const (
	Blue  = 0x5865f2
	Green = 0x2dcc70
	Red   = 0xe74d3b
)

const (
	MaxCoursesPerPage = 30
)

var EmojiMap = map[string]string{
	"spring": ":herb:",
	"summer": ":sunny:",
	"fall":   ":fallen_leaf:",
	"winter": ":snowflake:",
}

func LastOpenDisplayString(lastOpen int64) string {
	if lastOpen == 0 {
		return "`Currently`\n"
	} else if lastOpen == -1 {
		return "`Unknown`\n"
	} else {
		return fmt.Sprintf("<t:%d:R>\n", lastOpen)
	}
}
