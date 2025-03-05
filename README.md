# calculator-on-go

Этот проект представляет собой калькулятор, реализованный на языке Go. Система для параллельного вычисления сложных арифметических выражений с использованием оркестратора и агента.

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
curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d '{"expression": "(3 +2)*5 "}'

#Ответ(200) :
{"expressions":[{"id":"10c77400-cf81-48f9-ac58-e3c6574090dd","status":"completed","result":25,"input":"(3 +2)*5 "}]}
```

2. Ошибка деления на ноль
```bash
# Отправка выражения
curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d '{"expression": "3 / 0 "}'

# Ответ(422):
{"expressions":[{"id":"e61b6fae-d03b-483b-9f2b-5b64221ca0e6","status":"error","error":"деление на ноль","input":"3 / 0"}]}
```

3. Выражение с буквами

```bash
# Отправка выражения
curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d '{"expression": "df"}'

# Ответ(422):
выражение неверно: разрешены только числа и ( ) + - * /
```
4. Выражение c лишними знаками действия, с не закрытой скобкой

```bash
# Отправка выражения
curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d '{"expression": "2++5"}'

# Ответ(422):
недостаточно операндов для операции +
```

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
├─ .evn
├─ agent
├─ go.mod
├─ go.sum                
├─ orchestrator                             
├─ README.md
```