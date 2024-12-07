package main

import (
	"log"
	"net/http/cookiejar"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	version  string
	cookie   *cookiejar.Jar
}

var app = &application{
	version: "v0.1.0",
}

func main() {
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.cookie, _ = cookiejar.New(nil)

	app.parseArgs()
}
