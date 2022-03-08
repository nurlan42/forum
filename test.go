package main

import (
	"log"
	"os"
)

func main() {
	bold := "\033[1m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	reset := "\033[0m"
	InfoLogger := log.New(os.Stdout, bold+colorGreen+"INFO:\t "+reset, log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, bold+colorRed+"ERROR: \t"+reset, log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger.Println("Hello info")
	ErrorLogger.Println("Hello error")
}
