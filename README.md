## Стек проекта

СУБД Postresql, Gin web framework, gorm ORM 

Сервер и база данных запускаются в докере

Сервер запускается в окружении разработчика

## Запуск проекта
Запускать из папки где лежит docker-compose.yml
```bash
docker compose up
```
Если нужно запустить отдельно используйте флаг -d без отображения логов в лайве

Логи докера можно просмотреть с помощью команды:
```bash
docker compose logs -ft
```


В конфиге docker-compose также есть сервис который я использовал для теста в качестве внешнего API
Его можно закоментить и запустить без него.
В этом случае не забудьте указать базовый URL на внешний API машин в .env файле

В моем случаем он выглядит так.
```
PORT=4000

DB_URL="host=e_m_postgres_db user=postgres password=postgres dbname=postgres port=4001 sslmode=disable"
DB_USER=postgres
DB_PASS=postgres
DB_NAME=postgres
DB_HOST=e_m_postgres_db
DB_PORT=4001
DB_SSL=disable

GO_ENV=DEV

EXTERNAL_CAR_API_URL=http://json-server:3000
```

## Тестирование эндпоинтов
Открыть коллекцию эндпоинтов в Postman 
[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/18919361-3345e798-59ed-426d-92cc-993e980e589f?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D18919361-3345e798-59ed-426d-92cc-993e980e589f%26entityType%3Dcollection%26workspaceId%3Dae74495c-6346-4690-9da5-9e459609642c)

Ссылка на эндпоинты в сваггере при запущенном проекте [Тык](http://localhost:4000/swagger/index.html#/)

При добавлении новых эндпоинтов и их описания нужно запускать команду:
```bash
swag init
```
Данная команда может быть выполнена при запущенном поекте. CompileDaemon перебилдит проект после новой инициализации файлов сваггера.

## Логи
Логи Gin записываются в backend_info.log

Логи Gorm записываются в gorm.log