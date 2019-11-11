# Тестовое задание для [Clickadu](https://clickadu.com/)

Текст самого задания: [посмотреть](/task.md)

# Установка
1. `mkdir -p $GOPATH/src/github.com/Madredix/clickadu && cd $GOPATH/src/github.com/Madredix/clickadu`
2. `git clone https://github.com/Madredix/clickadu ./`
3. `go get github.com/go-swagger/go-swagger` при отсутствии go-swagger
4. `make build`

# Запуск
* `./app -c=./config.json.example`
* `./app --help` - посмотреть доступные команды

# Запуск тестов
* `make test-unit`

# Запуск в Docker
* `docker build -t clickadu . && docker run -it clickadu`

# Примеры запросов
* `curl -X POST -H "Content-Type: application/json" -d '[{"url": "http://yandex.ru", "number_of_requests": 3},{"url": "http://google.com", "number_of_requests": 3}]' http://localhost:2000/add`
* `curl -X GET -H "Content-Type: application/json" http://localhost:2000/status`

# Комментарии
1. Для ускорения не сделал загрузку/сохранение статистики, теряется при перезапуске
2. Так же для ускорения оставил дефолтный веб-сервер, в нем нет ни логирования, ни метрик
3. Не стал изменять swagger из задания, я бы поменял:
    1. для поля `url` нужно прописать `"format": "uri"` что бы swagger валидировал входящие url
    2. для схемы `/add` прописать `"minLength": 1` так как нет смысла добавлять пустую задачу 
    3. вместо 405 _(Method Not Allowed, а не как написано Invalid input)_ для `/add` сделать 400 (если не верные данные, в текущей реализации не валидируются) + 500 если не удалось добавить в очередь (в текущей реализации нельзя не добавить, но предусмотреть нужно)
    4. прописать `"x-omitempty": false` для всех полей статистики, что бы api всегда отображал все поля из схемы (нужно, например, при тестировании схемы)
4. Как превратить из тестового в реальный проект:
    1. Сохранение/загрузка статистики при запуске/остановке/периодически. Например, в consul - это позволит запускать несколько инстансов приложения за балансером и видеть по всем статистику отдельно/суммарно
    2. Очередь вынести в отдельное приложение, например, kafka
    3. Обработчик очереди отдельным инстансом с отдельной консольной командой запуска
5. Тесты для Newman не успел написать