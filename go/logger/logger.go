package logger

import (
	"log"
	"os"
)

// Debug is the debug logger that writes to debug.log.
// It is initialized in main() so gameloop and other packages can use it.
var Debug *log.Logger

// Init opens (or creates) debug.log in append mode and sets up Debug.
func Init() error {
	f, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	Debug = log.New(f, "[DEBUG] ", log.Ltime)
	return nil
}
