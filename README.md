## Сервис задач согласования и отправки писем 

* реализует REST API для создания задачи на согласование и рассылки уникальных ссылок-решений участникам. 
* Список участников указывается явно (список email) для каждой задачи. 
* Аутентификация обращений на REST API валидируется на сервисе аутентификации посредством GRPC-вызовов (отключено в связи для возможности независимого запуска сервиса). 
* REST API должно реализовывать CRUDL для задач согласования. Операции U и D позволены только автору задачи. 
* Каждому участнику при создании/обновлении задачи высылается письмо с двумя уникальными ссылками - "согласовано" и "не согласовано". Сначала отправляется письмо первому согласующему, ожидается его реакция, затем следующему, ожидается его реакция, и т.д. до последнего согласующего. API должно иметь методы для обработки "нажатий" на высланные ссылки и регистрации соответствующей реакции согласующего. 
* Если была нажата хотя бы одна ссылка "не согласовано", задача считается в целом не согласованной, и всем участникам рассылается уведомление об окончании согласования с негативным результатом, и письма со ссылками следующим согласующим по этой задаче уже не отправляются. 
* Сервис формирует и отправляет в kafka события создания задач, отправки писем, нажатия на ссылки.

---

## Запуск с использованием Docker

### Настройка проекта

Создайте `.env` файл в корне репозитория:

```bash
cp .env.example .env
```

Внесите при необходимости корректировки в переменные окружения.

### Сборка образов и запуск контейнеров

В корне репозитория выполните команду:

```bash
docker-compose up --build
```

### Остановка контейнеров

Для остановки контейнеров выполните команду:

```bash
docker-compose stop
```

---

Документация - http://localhost:3000/swagger/index.html

Запросы следует отправлять на http://localhost:3000/
