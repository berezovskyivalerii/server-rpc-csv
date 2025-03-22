#  gRPC CSV Product Service

## server-rpc-csv — gRPC-сервер

### Возможности:
- `Fetch(URL)` — загружает CSV-файл по ссылке, сохраняет данные о продуктах в MongoDB.
- `List(page, size, sort)` — возвращает список продуктов с пагинацией и сортировкой.

### Используемый стек:

- Go 1.13+
- MongoDB
- gRPC + Protobuf
- `encoding/csv`

### Как запустить сервер

#### 1. Установить зависимости

```bash
go mod tidy
```
#### Установить переменные окружения 
```bash
export DB_URI=mongodb://localhost:27017
export DB_DATABASE=product
export DB_USERNAME=admin
export DB_PASSWORD=cheapshots
export SERVER_PORT=9000
```

#### 3. Запустить сервер
```bash
go run cmd/main.go
```

#### Пример CSV-файла:
```csv
Name,Price
milk,27
bread,18
cheese,50
```

#### Приложение рабоатет в связки с [клиентом](https://github.com/berezovskyivalerii/client-rpc-csv.git) и [CRUD приложением](https://github.com/berezovskyivalerii/csv-rest-app.git)

#### Выполняемая задача
```txt
Требуется написать gRPC-сервер на языке GoLang (1.13+), с постоянным хранилищем
MongoDB, реализующий 2 метода:

- Fetch(URL) - запросить внешний CSV-файл со списком продуктов по внешнему адресу.
CSV-файл имеет вид PRODUCT NAME;PRICE. Последняя цена каждого продукта должна
быть сохранена в базе с датой запроса. Также нужно сохранять количество изменений
цены продукта.

- List(<paging params>, <sorting params>) - получить постраничный список продуктов с их
ценами, количеством изменений цены и датами их последнего обновления.
Предусмотреть все варианты сортировки для реализации интерфейса в виде
бесконечного скролла.
```
