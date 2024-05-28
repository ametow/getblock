# GetBlock.io

## Запуск

Для подключения к api getblock.io нам нужно указать api ключ в `makefile` в переменную `$API_KEY` или присвоить к переменному окружения `GETBLOCK_KEY`

### Консольная версия

```sh
make cli
make run-cli
```

### Веб api версия

```sh
make service
make run-service
```

Для веб версии дергаем ручку по умолчанию: http://localhost:8080/api/get


## Docker

```sh
make container-cli
```

or

```sh
make container-service
```
