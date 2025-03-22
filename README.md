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
