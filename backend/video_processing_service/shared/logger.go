// File: shared/logger.go
package shared

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func InitLogger() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
