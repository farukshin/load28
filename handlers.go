package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (app *application) getVersion() {
	fmt.Printf("loader28 %s\n", app.version)
}

func (app *application) parseArgs() error {

	if len(os.Args) < 1 || isArgs("--help") || isArgs("-h") {
		app.help_home()
	} else if isArgs("--version") || isArgs("-v") {
		app.getVersion()
	} else {
		app.run()
	}
	return nil
}

func initArgs(str string, env string) (string, error) {

	val, errh := getArgs(str)
	if errh != nil {
		val = os.Getenv(env)
		if val == "" {
			return "", fmt.Errorf("Не задан параметр %s. Задайте параметр %s в команде запуска или иницализируйте переменную окружения %s", str, str, env)
		}
	}
	return val, nil
}

func (app *application) run() {

	user, err := initArgs("--user", "LOADER28_USER")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	password, err := initArgs("--password", "LOADER28_PASSWORD")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	debug, _ := initArgs("--debug", "LOADER28_PASSWORD")

	command, err := getCommand()
	if err != nil {
		app.help_home()
		return
	}
	if command == "list" {
		app.list(user, password, debug)
	} else if command == "get" {
		app.get(user, password, debug)
	} else {
		app.help_home()
	}
}

func (app *application) list(login string, password string, debug string) {

	client := &http.Client{
		Jar: app.cookie,
	}

	type loginP struct {
		Login       string `json:"login"`
		Password    string `json:"password"`
		ServiceNick string `json:"serviceNick"`
	}

	type ticket struct {
		Ticket string `json:"ticket"`
	}

	ticketUrl := "https://login.1c.ru/rest/public/ticket/get"
	req, _ := http.NewRequest("POST", ticketUrl, nil)
	postBody, err := json.Marshal(
		loginP{login, password, "Platform83"})
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	req.SetBasicAuth(login, password)
	req.Header.Set("Content-Type", "application/json")
	buf := bytes.NewBuffer(postBody)
	req, err = http.NewRequest("POST", ticketUrl, buf)

	resp, _ := client.Do(req)

	var ticketData ticket
	b, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.errorLog.Println(err)
	}
	json.Unmarshal(b, &ticketData)

	aa := fmt.Sprintf("https://login.1c.ru//ticket/auth?token=%s", ticketData.Ticket)
	println(aa)
	return

}

func (app *application) get(user string, password string, debug string) {
	fmt.Println("get")
}

func getCommand() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("Не задана команда")
	}
	return os.Args[1], nil
}

func isArgs(a1 string) bool {

	_, err := getArgs(a1)
	return err == nil

}

func getArgs(a1 string) (string, error) {

	for _, s := range os.Args[1:] {
		if s == a1 {
			return "", nil
		}
		for i := 0; i < len(s); i++ {
			if s[i] == '=' && i > 0 {
				v := s[:i]
				if v == a1 {
					return s[i+1:], nil
				}
			}
		}

	}
	return "", errors.New("Не найдено флага " + a1)
}
func (app *application) help_home() {
	fmt.Println(`Приложение: loader28\n
    Загрузка дистрибутивов с сайта releases.1c.ru
	
Строка запуска: parser1c [КОМАНДА] [ОПЦИИ]
КОМАНДА:
    list - вывод списка доступных дистрибутивов
    get - загрузка дистрибутива

ОПЦИИ:
    -h --help - вызов справки
    -v --version - версия приложения
    --user - пользователь PostgreSQL (либо env PG_USER)
    --password - пароль PostgreSQL (либо env PG_PASSWORD)

Пример запуска:
    ./loader28 list --user=user1c --password=pass1c
    ./loader28 get --user=user1c --password=pass1c`)
}
