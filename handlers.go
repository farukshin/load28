package main

import (
	"errors"
	"fmt"
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
	println(
		"user:", user,
		"password:", password,
		"debug:", debug,
	)
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
