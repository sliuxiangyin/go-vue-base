// logger/logger.go
package logger

import "log"

func NewLogger() *log.Logger {
	return log.Default()
}
