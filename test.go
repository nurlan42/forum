package main

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func main() {
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	reset := "\033[0m"
	bold := "\033[1m"
	InfoLogger = log.New(os.Stdout, bold+colorGreen+"INFO: "+reset, log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, bold+colorRed+"WARNING: "+reset, log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger.Println("Hi")
	WarningLogger.Println("Hello")

}
