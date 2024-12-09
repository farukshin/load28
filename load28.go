package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type buid struct {
	name     string `json:"name"`
	url      string `json:"url"`
	fileName string `json:"fileName"`
}

type release struct {
	name    string `json:"name"`
	version string `json:"version"`
}

type soft struct {
	name string `json:"name"`
	nick string `json:"nick"`
}

func (app *application) get() {

	err := app.auth()
	if err != nil {
		app.errorLog.Println(err)
	}

	if app.soft == "" || app.release == "" {
		app.help_home()
		return
	}
	buids, err := app.getBuids()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	err = app.downloads(buids[0])
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) downloads(b buid) error {

	req, err := http.NewRequest("GET", app.urlReleases+b.url, nil)
	if err != nil {
		return err
	}
	resp, err := app.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	out, err := os.Create(b.fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	return nil
}

func (app *application) list() {

	err := app.auth()
	if err != nil {
		app.errorLog.Println(err)
	}

	if app.soft == "" {
		res, err := app.getSofts()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		for _, v := range res {
			fmt.Printf("%s=%s\n", v.nick, v.name)
		}
	} else if app.release == "" {
		rel, err := app.getReleases()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		for _, v := range rel {
			fmt.Printf("%s=%s\n", v.name, v.version)
		}
	} else if app.release != "" {
		buids, err := app.getBuids()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		for _, v := range buids {
			fmt.Printf("%s=%s\n", v.name, v.url)
		}
	}
}

func (app *application) getBuids() ([]buid, error) {
	res := make([]buid, 0)
	url := fmt.Sprintf("%s/version_files?nick=%s&ver=%s", app.urlReleases, app.soft, app.release)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, err
	}
	resp, err := app.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return res, err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	bb := string(strings.Replace(string(b), "\n", "", -1))
	bb = string(strings.Replace(bb, "  ", "", -1))
	bb = string(strings.Replace(bb, "\t", "", -1))
	bb = string(strings.Replace(bb, "<a href=", "\n<a href=", -1))

	reg := fmt.Sprintf(`<a href="(/version_file\?nick=%s&ver=%s&path=.*)">(.*)</a>`, app.soft, app.release)
	r := regexp.MustCompile(reg)
	matches := r.FindAllStringSubmatch(bb, -1)

	for _, v := range matches {
		ind := strings.LastIndex(v[1], "%5c")
		fileName := v[1]
		if ind != -1 {
			fileName = v[1][ind+3:]
		}
		res = append(res, buid{name: v[2], url: v[1], fileName: fileName})
	}

	return res, nil
}

func (app *application) getReleases() ([]release, error) {
	url := fmt.Sprintf("%s/project/%s?allUpdates=true", app.urlReleases, app.soft)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := app.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	bb := string(strings.Replace(string(b), "\n", "", -1))
	bb = string(strings.Replace(bb, "  ", "", -1))
	bb = string(strings.Replace(bb, "\t", "", -1))

	reg := fmt.Sprintf(`<a href="/version_files\?nick=%s&ver=(\S*)">(\S*)</a>`, app.soft)
	r := regexp.MustCompile(reg)
	matches := r.FindAllStringSubmatch(bb, -1)
	res := make([]release, 0)
	for _, v := range matches {
		fmt.Printf("%s=%s\n", v[1], v[2])
		res = append(res, release{name: v[2], version: v[1]})
	}
	return res, nil
}

func (app *application) getSofts() ([]soft, error) {
	url := app.urlReleases
	if app.hideUnavailablePrograms {
		url += "?hideUnavailablePrograms=true"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := app.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadAll(resp.Body)

	r := regexp.MustCompile(`<a href="/project/(.*)">(.*)</a>`)
	matches := r.FindAllStringSubmatch(string(b), -1)
	res := make([]soft, 0)
	for _, v := range matches {
		fmt.Printf("%s=%s\n", v[1], v[2])
		res = append(res, soft{name: v[2], nick: v[1]})
	}
	return res, nil
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
