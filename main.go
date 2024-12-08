package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
)

type application struct {
	errorLog                *log.Logger  `json:"errorLog"`
	infoLog                 *log.Logger  `json:"infoLog"`
	version                 string       `json:"version"`
	client                  *http.Client `json:"client"`
	login                   string       `json:"login"`
	pwd                     string       `json:"pwd"`
	soft                    string       `json:"soft"`
	release                 string       `json:"release"`
	hideUnavailablePrograms bool         `json:"hideUnavailablePrograms"`
}

var app = &application{
	version: "v0.1.0",
}

func main() {
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	cookie, _ := cookiejar.New(nil)
	app.client = &http.Client{
		Jar: cookie,
	}
	app.parseArgs()
}
