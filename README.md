# Тестовое задание на позицию Go разработчика
[![Build Status](https://travis-ci.org/dimuska139/cacher.svg?branch=master)](https://travis-ci.org/dimuska139/cacjer)
[![codecov](https://codecov.io/gh/dimuska139/cacher/branch/master/graph/badge.svg)](https://codecov.io/gh/dimuska139/cacher)

## Задача
Написать простой GRPC сервис с командами: **get**, **set** и **delete**.
Хранилище должно быть спрятано за интерфейсом.

Реализуйте интерфейс обоими путями:

- Memcached сервер с самописной библиотекой и тремя этими же командами. [Memcached protocol](https://github.com/memcached/memcached/blob/master/doc/protocol.txt)
- Хранилище внутри памяти приложения

Оформите в виде git-репозитория и покройте тестами.

Продвинутый уровень: реализуйте пулл коннектов к memcached.

## Реализация

Самописная библиотека для работы с Memcache находится в директории
`libs/memcache`. Пулл коннектов в ней реализован. Хранилища находятся в директории
`internal/cache`. Интерфейс к ним находится там, где они используются - то есть
в `internal/api/grpc/cache_server.go`. Выбор типа используемого хранилища
осуществляется с помощью переменной в конфигурационном файле - `storage`. Proto-файлы находятся
тут: `internal/api/grpc/proto`.

## Запуск
1. Скопировать файл `config.yml.dist` (это шаблон) в `config.yml`
2. Запустить docker-compose: `sudo docker-compose up -d`
3. Запустить сервис: `make run`

## Команды
* `make run` - запуск сервиса
* `make test` - запуск тестов
* `make grpc_server` - генерация proto