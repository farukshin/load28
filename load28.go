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

	if app.soft == "" {
		app.errorLog.Println("Не задан параметр --soft")
		return
	}
	if app.release == "" {
		rel, err := app.getReleases()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		app.release = rel[0].version
		if app.debug {
			app.infoLog.Printf("Не задан параметр --release, инициализирую параметр релиз значением %s\n", app.release)
		}
	}

	buids, err := app.getBuids()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	if len(buids) == 0 {
		app.errorLog.Println("Не найдено билдов")
		return
	} else if len(buids) > 1 && app.debug {
		app.infoLog.Printf("Внимание! Получил более одного билда для загрузки, загружаю только первый билд '%s'\n", buids[0].name)
	}
	err = app.downloads(buids[0])
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Println("success")
}

func (app *application) getResp(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err = app.client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (app *application) downloads(b buid) error {

	if app.debug {
		app.infoLog.Printf("Начинаю загрузку файла '%s'\n", b.fileName)
	}
	url := app.urlReleases + b.url
	resp, err := app.getResp(url)
	if err != nil {
		return err
	}
	bb, _ := ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile(`<a href="(https://.*)">(Скачать дистрибутив)</a>`)
	matches := r.FindAllStringSubmatch(string(bb), -1)
	res := make([]string, 0)
	for _, v := range matches {
		res = append(res, v[1])
	}
	if len(res) == 0 {
		return fmt.Errorf("Не найдена ссылка на скачивание файла '%s'", b.fileName)
	}

	resp, err = app.getResp(res[0])
	if err != nil {
		return err
	}
	if app.debug {
		app.infoLog.Printf("Сохраняю файл на диске '%s'\n", b.fileName)
	}
	out, err := os.Create(b.fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	if app.debug {
		app.infoLog.Printf("Загрузка файла успешно завершена '%s'\n", b.fileName)
	}
	return nil
}

func (app *application) list() {

	err := app.auth()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	if app.soft == "" {
		res, err := app.getSofts()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		if !app.debug {
			for _, v := range res {
				fmt.Printf("%s=%s\n", v.nick, v.name)
			}
		}
	} else if app.release == "" {
		rel, err := app.getReleases()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		if !app.debug {
			for _, v := range rel {
				fmt.Printf("%s=%s\n", v.name, v.version)
			}
		}
	} else if app.release != "" {
		buids, err := app.getBuids()
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		if !app.debug {
			for _, v := range buids {
				fmt.Printf("%s=%s\n", v.name, v.url)
			}
		}
	}
}

func (app *application) getBuids() ([]buid, error) {

	if app.debug {
		app.infoLog.Printf("Начинаю загрузку билдов программного обеспечения '%s' версии релиза '%s'\n", app.soft, app.release)
	}
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
	str := ""
	for _, v := range matches {
		ind := strings.LastIndex(v[1], "%5c")
		fileName := v[1]
		if ind != -1 {
			fileName = v[1][ind+3:]
		}
		str += fmt.Sprintf("\t%s\n", v[2])
		res = append(res, buid{name: v[2], url: v[1], fileName: fileName})
	}
	if app.debug {
		app.infoLog.Printf("Получил список билдов программного обеспечения '%s' (%d шт):\n%s", app.soft, len(matches), str)
	}
	if app.filter == "" {
		return res, nil
	}
	resFilter, err := filterBuilds(res)
	return resFilter, nil
}

func filterBuilds(buids []buid) ([]buid, error) {

	if app.debug {
		app.infoLog.Printf("Начинаю фильтрацию по регулярному выражению '%s' билдов программного обеспечения '%s' версии релиза '%s'\n", app.filter, app.soft, app.release)
	}
	res := make([]buid, 0)
	str := ""
	for _, v := range buids {
		match, _ := regexp.MatchString(app.filter, v.name)
		if match {
			res = append(res, v)
			str += fmt.Sprintf("\t%s=%s\n", v.name, v.url)
		}
	}
	if app.debug {
		app.infoLog.Printf("Получил отфильтрованный список билдов программного обеспечения '%s' (%d шт):\n%s", app.soft, len(res), str)
	}
	return res, nil
}

func (app *application) getReleases() ([]release, error) {

	if app.debug {
		app.infoLog.Printf("Начинаю загрузку релизов программного обеспечения '%s'\n", app.soft)
	}
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
	str := ""
	for _, v := range matches {
		str += fmt.Sprintf("\t%s=%s\n", v[1], v[2])
		res = append(res, release{name: v[2], version: v[1]})
	}
	if app.debug {
		app.infoLog.Printf("Получил список доступных релизов программного обеспечения '%s' (%d шт.):\n%s", app.soft, len(matches), str)
	}
	if app.filter == "" {
		return res, nil
	}
	resFilter, err := filterReleases(res)
	return resFilter, nil
}

func filterReleases(releases []release) ([]release, error) {

	if app.debug {
		app.infoLog.Printf("Начинаю фильтрацию по регулярному выражению '%s' программного обеспечения '%s' версии релиза '%s'\n", app.filter, app.soft, app.release)
	}
	res := make([]release, 0)
	str := ""
	for _, v := range releases {
		match, _ := regexp.MatchString(app.filter, v.name)
		if match {
			res = append(res, v)
			str += fmt.Sprintf("\t%s=%s\n", v.name, v.version)
		}
	}
	if app.debug {
		app.infoLog.Printf("Получил отфильтрованный список релизов программного обеспечения '%s' (%d шт):\n%s", app.soft, len(res), str)
	}
	return res, nil
}

func (app *application) getSofts() ([]soft, error) {

	if app.debug {
		app.infoLog.Printf("Начинаю загрузку списка доступного программного обеспечения с сайта %s\n", app.urlReleases)
	}
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
	po := ""
	for _, v := range matches {
		po += fmt.Sprintf("\t%s=%s\n", v[1], v[2])
		res = append(res, soft{name: v[2], nick: v[1]})
	}
	if app.debug {
		app.infoLog.Printf("Получил список доступного программного обеспечения (%d шт.):\n%s", len(matches), po)
	}
	if app.filter == "" {
		return res, nil
	}
	resFilter, err := filterSofts(res)
	return resFilter, nil
}

func filterSofts(softs []soft) ([]soft, error) {

	if app.debug {
		app.infoLog.Printf("Начинаю фильтрацию по регулярному выражению '%s' доступного программного обеспечения\n", app.filter)
	}
	res := make([]soft, 0)
	str := ""
	for _, v := range softs {
		match, _ := regexp.MatchString(app.filter, v.name)
		if match {
			res = append(res, v)
			str += fmt.Sprintf("\t%s=%s\n", v.nick, v.name)
		}
	}
	if app.debug {
		app.infoLog.Printf("Получил отфильтрованный доступного программного обеспечения(%d шт):\n%s", len(res), str)
	}
	return res, nil
}

func (app *application) auth() error {

	if app.debug {
		app.infoLog.Printf("Запускаю прохождение авторизации на сайте %s\n", app.urlLogin)
	}
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
	if app.debug {
		app.infoLog.Printf("Авторизация успешно пройдена\n")
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
