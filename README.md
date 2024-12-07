# Loader28

Загрузка дистрибутивов с сайта releases.1c.ru

* [Установка](#install)
* * [Установка из исходников](#installSource)
* * [Установка из релизов](#installRelease)
* * [Установка из docker образа](#installDocker)
* [Использование](#usage)
* * [Запуск из консоли](#runcli)
* * [Запуск в docker образе](#rundocker)
* [Нашли ошибку?](#err)
* [Лицензия](#lic)


<a name="install"></a> 

## Установка

<a name="installSource"></a> 

### Установка из исходников

```
git clone https://github.com/farukshin/load28.git
cd load28
go build .
./load28 --version
```

<a name="installRelease"></a> 

### Установка из релизов

1. Получить версию [последнего релиза](https://github.com/farukshin/load28/releases).

``` bash
VERSION=$(curl -s "https://api.github.com/repos/farukshin/load28/releases/latest" | jq -r '.tag_name')
```
Или установить необходимую версию релиза:

``` bash
VERSION=vX.Y.Z
```

2. Загрузка релиза

``` bash
OS=Linux       # or Darwin, Windows
ARCH=x86_64    # or arm64, x86_64, armv6, i386, s390x
FILE=load28_${OS}_${ARCH}.tar.gz
curl -sL "https://github.com/farukshin/load28/releases/download/${VERSION}/${FILE}" > ${FILE}
```

3. Проверка контрольной суммы

``` bash
curl -sL https://github.com/farukshin/load28/releases/download/${VERSION}/load28_checksums.txt > load28_checksums.txt
shasum --check --ignore-missing ./load28_checksums.txt
```

4. Распаковать утилиту

``` bash
tar -zxvf ${FILE} load28
./load28 --version
```

<a name="installDocker"></a> 

### Установка из docker образа

`load28` можно запустить из docker образа. Сам образ можно скачать из docker hub'a

```
docker push farukshin/load28
```
или собрать локально

```
git clone https://github.com/farukshin/load28.git
cd load28
docker build -t farukshin/load28 .
```

Образ `farukshin/load28` создан на базе `scratch`, поэтому итоговый размер образа 5MB
```
docker images | grep "farukshin/load28"
> farukshin/load28   latest    8492ada2a558   7 minutes ago   5.73MB
```

<a name="usage"></a> 

## Использование

<a name="runcli"></a> 

### Запуск из консоли

```
Строка запуска: load28 [КОМАНДА] [ОПЦИИ]
КОМАНДА:
    list - вывод списка доступных дистрибутивов
    get - загрузка дистрибутива

ОПЦИИ:
    -h --help - вызов справки
    -v --version - версия приложения
    --login - пользователь портала releases.1c.ru (либо env LOAD28_USER)
    --pwd - пароль пользователя портала releases.1c.ru (либо env LOAD28_PWD)
```

Вывод списка доступных дистрибутивов:

```
export LOAD28_USER=myuser
export LOAD28_PWD=mypassword
./load28 list | head -n 10
```

Результат:
```
DevelopmentTools10=1C:Enterprise Development Tools
Executor=1C:Исполнитель
Analytics=1С:Аналитика
Conversion=1С:Конвертация данных 2.0
Conversion30=1С:Конвертация данных 3
Translator=1С:Переводчик, редакция 2.1
CollaborationSystem=1С:Сервер взаимодействия
STest=1С:Сценарное тестирование 8
Tester=1С:Тестировщик
esb=1С:Шина
```

<a name="rundocker"></a> 

### Запуск в docker образе

```
docker run --rm \
    -e LOAD28_USER='myuser' \
    -e LOAD28_PWD='mypassword' \
    farukshin/load28:latest list
```

<a name="err"></a> 

## Нашли ошибку?

Если при использовании `load28` нашли ошибку - создайте [новый issues](https://github.com/farukshin/load28/issues/new). 

<a name="lic"></a> 

## Лицензия

`load28` выпускается под лицензией MIT. Подробнее [LICENSE.md](https://github.com/farukshin/load28/blob/main/LICENSE.md)
