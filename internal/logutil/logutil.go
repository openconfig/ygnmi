package logutil

import (
	"bufio"
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
	scanner := bufio.NewScanner(strings.NewReader(message))
	for scanner.Scan() {
		line := scanner.Text()
		length := len(line)
		for startIndex := 0; startIndex < length; startIndex += logLimit {
			endIndex := startIndex + logLimit
			if endIndex > length {
				endIndex = length
			}
			log.V(logLevel).Info(line[startIndex:endIndex])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Warningf("ygnmi/logutil: error while scanning log input: %v", err)
	}
}
