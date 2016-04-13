// +build develop

package log

import (
	. "log"
	"os"
)

var (
	format int = Ldate | Ltime | Llongfile

	Info  *Logger = New(os.Stdout, "INFO: ", format)
	Debug *Logger = New(os.Stdout, "DEBUG: ", format)
	Error *Logger = New(os.Stderr, "ERROR: ", format)
)
