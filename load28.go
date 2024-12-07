package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func (app *application) list() {

	type loginP struct {
		Login       string `json:"login"`
		Password    string `json:"password"`
		ServiceNick string `json:"serviceNick"`
	}

	err := app.auth()
	if err != nil {
		app.errorLog.Println(err)
	}

	if app.soft == "" {
		app.releases_home()
	}
}

func (app *application) releases_home() {
	url := "https://releases.1c.ru/"
	if app.hideUnavailablePrograms {
		url += "?hideUnavailablePrograms=true"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.errorLog.Println(err)
	}
	resp, err := app.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		app.errorLog.Println(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)

	r := regexp.MustCompile(`<a href="/project/(.*)">(.*)</a>`)
	matches := r.FindAllStringSubmatch(string(b), -1)
	for _, v := range matches {
		fmt.Printf("%s=%s\n", v[1], v[2])
	}
}

func (app *application) auth() error {

	token, err := app.getToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://login.1c.ru/ticket/auth?token=%s", token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(app.login, app.pwd)
	resp, err := app.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func (app *application) getToken() (string, error) {

	type tokenReq struct {
		Login       string `json:"login"`
		Password    string `json:"password"`
		ServiceNick string `json:"serviceNick"`
	}

	type ticket struct {
		Ticket string `json:"ticket"`
	}

	ticketUrl := "https://login.1c.ru/rest/public/ticket/get"
	postBody, err := json.Marshal(tokenReq{app.login, app.pwd, "https://releases.1c.ru"})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ticketUrl, bytes.NewReader(postBody))
	req.SetBasicAuth(app.login, app.pwd)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.client.Do(req)
	defer resp.Body.Close()

	var ticketData ticket
	b, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.errorLog.Println(err)
	}
	err = json.Unmarshal(b, &ticketData)

	return ticketData.Ticket, err
}
