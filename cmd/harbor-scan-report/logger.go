package main

import (
	"log"
	"os"
)

type Loggers struct {
	Error   *log.Logger
	Warning *log.Logger
	Info    *log.Logger
	Debug   *log.Logger
}

var Log Loggers

func initLoggingSystem() {
	Log := Loggers{
		Error:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warning: log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Debug:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	Log.Info.Println("Starting the application...")
	Log.Info.Println("Logging system initialized")
}
