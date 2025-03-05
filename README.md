# calculator-on-go

Этот проект представляет собой калькулятор, реализованный на языке Go. Он позволяет вычислять математические выражения, такие как сложение, умножение, деление, вычитание и возведение в степень.

## Возможности

- Поддержка сложения (`+`), вычитания (`-`), умножения (`*`), деления (`/`) и возведения в степень (`^`).
- Обработка математических выражений, включая вложенные скобки.
- Возвращает подробные сообщения об ошибках для некорректных выражений.

---

## Запуск проекта

Следуйте этим шагам, чтобы запустить проект:

1. Клонируйте репозиторий:
 
```bash
git clone https://github.com/Maraei/calculator-on-go.git
cd calculator-on-go
```

2. Запуск оркестратора:

```bash
go run ./cmd/orchestrator/main.go
```
3. Запуск агента:
Linux / macOS (bash):
```bash
export ORCHESTRATOR_URL=http://localhost:8080
go run ./cmd/agent/main.go
```
Windows (cmd.exe):
```bash
set ORCHESTRATOR_URL=http://localhost:8080
go run .\cmd\agent\main.go
```

## Примеры использования

1. Добавление выражения

```bash
curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d '{"expression": "3 + 4 * 2"}'
```
```
Ответ:
```json
{
  "id": "1"
}
```
2. Получение списка выражений
```bash
curl -X GET "http://localhost:8080/api/v1/expressions"
```

## Тестирование

Чтобы выполнить тесты, используйте команду:

```bash
go test ./...
```

## Примеры сценариев
1. Успешное вычисление
```bash
# Отправка выражения
curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d '{"expression": "3 * 2 +5 "}'

# Проверка статуса
curl http://localhost:8080/api/v1/expressions

# Ответ через 500 мс:
{
    "expression": {
        "id": "1",
        "status": "completed",
        "result": 11
    }
}
```

2. Ошибка деления на ноль
```bash
# Отправка выражения
curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d '{"expression": "3 / 0 "}'

# Проверка статуса
curl http://localhost:8080/api/v1/expressions

# Ответ через 500 мс:
{
    "expression": {
        "id": "2",
        "status": "error",
        "result": null
    }
}

## Структура проекта

```
calculator-on-go
├─ cmd
│  ├─ agent
│  │  ├─ main.go
│  ├─ orchestrator
│  │  ├─ main_test.go
│  │  ├─ main.go                       
├─ internal
│  ├─ agent
│  │  ├─ agent_test.go            
│  │  ├─ calculator.go 
│  │  ├─ task_fetcher.go
│  ├─ orchestrator
│  │  ├─ handler.go            
│  │  ├─ service.go 
│  │  ├─ task_manager.go                
├─ pkg
│  ├─ models       
│  │  ├─ expression.go            
│  ├─ utils       
│  │  ├─ config.go
├─ .evn
├─ agent
├─ go.mod
├─ go.sum                
├─ orchestrator                             
├─ README.md
```