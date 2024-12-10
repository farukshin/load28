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

	release, err := getArgs("--release")
	if err == nil && release != "" {
		app.release = release
	}

	filter, err := getArgs("--filter")
	if err == nil && filter != "" {
		app.filter = filter
	}

	debug, err := getArgs("--debug")
	if err == nil && (debug == "1" || debug == "true") {
		app.debug = true
	}

	app.urlReleases = "https://releases.1c.ru"
	app.urlLogin = "https://login.1c.ru"
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
					return trimkov(s[i+1:]), nil
				}
			}
		}

	}
	return "", errors.New("Не найдено флага " + a1)
}

func trimkov(s string) string {

	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	if len(s) >= 4 && s[0] == '\\' && s[1] == '"' && s[len(s)-1] == '"' && s[len(s)-2] == '\\' {
		return s[2 : len(s)-2]
	}
	return s
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
    --soft - наименование продукта (пример: Platform83)
    --release - версия продукта (пример: 8.3.14.1855)
    --filter - фильтр поиска (регулярное выражение)
    --debug - режим отладки (для включения укажите 1 или true)

Пример запуска:
    export LOAD28_USER=user1c
    export LOAD28_PWD=pass1c
    ./load28 list
    ./load28 get --soft=Platform83 --release=8.3.26.1498 --filter="Сервер.*ARM.*RPM.*Linux" --debug=1`)
}
