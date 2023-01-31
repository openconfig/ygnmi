package logutil

import (
	"strings"

	log "github.com/golang/glog"
)

const (
	// logLimit is used by logging helpers to avoid log truncation on some systems.
	logLimit = 15000
)

// LogByLine logs the message line-by-line, and splits lines as well given the
// log limit.
func LogByLine(logLevel log.Level, message string) {
	for _, line := range strings.Split(message, "\n") {
		length := len(line)
		for startIndex := 0; startIndex < length; startIndex += logLimit {
			endIndex := startIndex + logLimit
			if endIndex > length {
				endIndex = length
			}
			log.V(logLevel).Info(line[startIndex:endIndex])
		}
	}
}
