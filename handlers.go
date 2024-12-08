package main

import (
	"errors"
	"fmt"
	"os"
)

func (app *application) getVersion() {
	fmt.Printf("load28 %s\n", app.version)
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
			return "", fmt.Errorf("Не задан параметр %s. Задайте параметр %s в команде запуска или инициализируйте переменную окружения %s", str, str, env)
		}
	}
	return val, nil
}
func (app *application) init() {
	login, err := initArgs("--login", "LOAD28_USER")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.login = login
	pwd, err := initArgs("--password", "LOAD28_PWD")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.pwd = pwd

	hideUnavailablePrograms, err := getArgs("--hideUnavailablePrograms")
	if err == nil && hideUnavailablePrograms == "true" {
		app.hideUnavailablePrograms = true
	}

	soft, err := getArgs("--soft")
	if err == nil && soft != "" {
		app.soft = soft
	}
}

func (app *application) run() {

	app.init()

	command, err := getCommand()
	if err != nil {
		app.help_home()
		return
	}

	if command == "list" {
		app.list()
	} else if command == "get" {
		app.get()
	} else {
		app.help_home()
	}
}

func (app *application) get() {
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
	fmt.Println(`Приложение: load28\n
    Загрузка дистрибутивов с сайта releases.1c.ru
	
Строка запуска: load28 [КОМАНДА] [ОПЦИИ]
КОМАНДА:
    list - вывод списка доступных дистрибутивов
    get - загрузка дистрибутива

ОПЦИИ:
    -h --help - вызов справки
    -v --version - версия приложения
    --login - пользователь портала releases.1c.ru (либо env LOAD28_USER)
    --pwd - пароль пользователя портала releases.1c.ru (либо env LOAD28_PWD)

Пример запуска:
    ./load28 list --user=user1c --password=pass1c
    ./load28 get --user=user1c --password=pass1c`)
}
