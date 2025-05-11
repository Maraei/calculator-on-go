# calculator-on-go

Этот проект представляет распределённый калькулятор, реализованный на Go. Он поддерживает параллельное выполнение сложных арифметических выражений с помощью оркестратора и агентов. Микросервисная архитектура использует gRPC и HTTP.
## Возможности

- Поддержка сложения (`+`), вычитания (`-`), умножения (`*`), деления (`/`) и возведения в степень (`^`).
- Обработка математических выражений, включая вложенные скобки.
- Возвращает подробные сообщения об ошибках для некорректных выражений.

---

* **Orchestrator** - принимает выражения, разбивает их на задачи и распределяет между агентами. Отвечает за очередь задач и их обработку.
* **Agent** - получает задачи от Orchestrator, выполняет вычисления и отправляет результат обратно.
* **Auth** - управляет регистрацией, входом и валидацией токенов JWT. Защищает доступ к API.
* **Gateway** - HTTP-интерфейс поверх gRPC-сервисов. Переводит внешние HTTP-запросы в gRPC-вызовы и обратно. Упрощает клиентам взаимодействие с системой.

## Запуск проекта

Следуйте этим шагам, чтобы запустить проект:

### 1. Клонируйте репозиторий:
 
```bash
git clone https://github.com/Maraei/calculator-on-go.git
cd calculator-on-go
```
Установите зависимости
```bash
go mod tidy
```
### 2. Запуск оркестратора:

```bash
go run ./cmd/orchestrator/main.go
```
### 3. Запуск агента:

```bash
go run ./cmd/agent/main.go
```
### 3. Запуск Auth:

```bash
go run cmd/auth/main.go
```
### 4. Запуск Gateway:

```bash
go run cmd/gateway/main.go
```
## Примеры использования

### 1. Регистрация
```bash
grpcurl -plaintext -d '{"username": "<ваш логин>", "password": "<ваш пароль>"}' \
localhost:50051 calculator.AuthCalculatorService/Register
```
```bash
#ответ
{
  "message": "User registered successfully"
}
```
### 2. Логин
```bash
grpcurl -plaintext -d '{"username": "<ваш логин>", "password": "<ваш пароль>"}' \
localhost:50051 calculator.AuthCalculatorService/Login
```
```bash
#ответ
{
  "token": <ваш_токен>
}
```
### 3. Добавление выражения

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <ваш_токен>" \
-d '{"expression": "2 + 2"}'
```
```bash
#ответ
{
    "id": <уникальный идентификатор выражения>
}
```
### 4. Получение выражения по его идентификатору
```bash
curl -X GET http://localhost:8080/api/v1/result/<id> -H "Authorization: Bearer <ваш_токен>"
```

## Тестирование

Чтобы выполнить тесты, используйте команду:

```bash
go test ./...
```

## Примеры сценариев
### 1. Успешное вычисление
```bash
# Отправка выражения
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY3MDk0MjUsImlhdCI6MTc0NjcwODgyNSwiaWQiOjE5LCJuYmYiOjE3NDY3MDg4MjV9.07TqVAmoR_DDg2IBXYq2EtR8mxfcxHMbUW9M5KlToxg" \
-d '{"expression": "2 + 2 * 3"}'

#Ответ(200) :
{
  "id": "10c77400-cf81-48f9-ac58-e3c6574090dd"
}
```
```bash
#Получение результата:
curl -X GET http://localhost:8080/api/v1/result/10c77400-cf81-48f9-ac58-e3c6574090dd\
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY3MDk0MjUsImlhdCI6MTc0NjcwODgyNSwiaWQiOjE5LCJuYmYiOjE3NDY3MDg4MjV9.07TqVAmoR_DDg2IBXYq2EtR8mxfcxHMbUW9M5KlToxg"
#Ответ (200):
{
  "result": 8,
  "status": "completed",
  "error": ":"
}

```

### 2. Ошибка деления на ноль
```bash
# Отправка выражения
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY3MDk0MjUsImlhdCI6MTc0NjcwODgyNSwiaWQiOjE5LCJuYmYiOjE3NDY3MDg4MjV9.07TqVAmoR_DDg2IBXYq2EtR8mxfcxHMbUW9M5KlToxg" \
-d '{"expression": "2 + 2 * 3"}'

# Ответ(422):
{
  "error": "деление на ноль"
}
```

## Структура проекта

```
calculator-on-go
├─ api
│  ├─ api
│  │  ├─ auth.pb.gw.go
│  │  ├─ auth.pb.go
│  │  ├─ auth_grpc.pb.go
│  ├─ auth.proto
├─ cmd
│  ├─ agent
│  │  ├─ main.go
│  ├─ auth
│  │  ├─ main.go
│  ├─ orchestrator
│  │  ├─ main_test.go
│  │  ├─ main.go                       
├─ internal
│  ├─ agent
│  │  ├─ agent_test.go            
│  │  ├─ calculator.go 
│  │  ├─ task_fetcher.go
│  ├─ auth
│  │  ├─ auth.go            
│  │  ├─ jwt.go 
│  │  ├─ middleware.go
│  │  ├─ server.go
│  ├─ orchestrator
│  │  ├─ handler.go            
│  │  ├─ service.go 
│  │  ├─ task_manager.go  
│  │  ├─ parser.go            
│  │  ├─ repository.go 
│  │  ├─ service.go 
│  │  ├─ service_test.go  
│  ├─ test
│  │  ├─ agent_task_test.go            
│  │  ├─ auth_integration_test.go               
├─ .evn
├─ agent
├─ go.mod
├─ go.sum                            
├─ README.md
```