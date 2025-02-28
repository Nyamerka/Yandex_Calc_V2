# Yandex-калькулятор V2.0

## Made by Бутер Бродский aka Nyamerka)

Данное api было написано в рамках продолжения прохождения курса разработки на языке Go от Яндекса как развитие предыдущего проекта.

### Как запустить?
Для начала клонируем проект к себе:
```zsh
git clone https://github.com/Nyamerka/Yandex_Calc_V2
```
> [!NOTE]
> Все команды запускаются из корневой папки проекта.

Запускаем `Docker`:
```zsh
docker-compose up
```
> [!TIP]
> Не забудьте установить приложение Docker на свой ПК :space_invader:!

При корректном запуске Вы получите:
```zsh
[+] Running 1/0
 ✔ Container yandex_calc_v20-calculator-1  Created                                                                                                                0.0s 
Attaching to yandex_calc_v20-calculator-1
```
> [!NOTE]
> Возможно, будет долговато :hourglass_flowing_sand:...

При необходимости отчистить кеш и подтянуть изменения проекта можно с помощью следующей команды:
```zsh
docker-compose build --no-cache
```

На компьютер будет установлен образ `Debian`, где будет функционировать весь проект. Все зависимости ставятся внутри контейнера. Контейнер тестируется перед запуском с помощью `go test`. Если Вы хотите протестировать проект самостоятельно, то можно использовать команду ниже, находясь в корневой папке проекта:
```zsh
go test ./internal/...
```
> [!NOTE]
> Проект проходит обязательное тестирование перед запуском. Если тестирование провалено - проект не запустится :smiling_imp:.

Результат тестирования должен выглядеть примерно так, однако при первом запуске приписка `(cached)` будет отсутствовать:
```zsh
ok      Yandex_Calc_V2.0/internal/app   (cached)
ok      Yandex_Calc_V2.0/internal/eval  (cached)
ok      Yandex_Calc_V2.0/internal/queue (cached)
ok      Yandex_Calc_V2.0/internal/stack (cached)
```

Для более подробного вывода можно добавить флаг `-v`:
```zsh
go test -v ./internal/...
```
> [!NOTE]
> Большое количество вывода, не пугайтесь :ghost:.

### Архитектура проекта

```mermaid
graph TD
    U[Пользователь] -->|POST /calculate| O[Оркестратор]
    U -->|GET /expressions| O
    O -->|GET /internal/task| A[Агент]
    subgraph COMPUTING_POWER
        W1[Worker 0]
        W2[Worker 1]
        W3[Worker 2]
        Wn[Worker ...]
    end
    A -->|COMPUTING_POWER| COMPUTING_POWER
    COMPUTING_POWER -->|RESULT| A
    A -->|POST /internal/task| O

    style U fill:#9acd32,stroke:#333,stroke-width:2px,color:#000
    style O fill:#add8e6,stroke:#333,stroke-width:2px,color:#000
    style A fill:#add8e6,stroke:#333,stroke-width:2px,color:#000
    style W1 fill:#ff4500,stroke:#333,stroke-width:2px,color:#000
    style W2 fill:#ff4500,stroke:#333,stroke-width:2px,color:#000
    style W3 fill:#ff4500,stroke:#333,stroke-width:2px,color:#000
    style Wn fill:#ff4500,stroke:#333,stroke-width:2px,color:#000

    linkStyle 0 stroke:#228b22,stroke-width:2px,color:#228b22
    linkStyle 1 stroke:#4682b4,stroke-width:2px,color:#4682b4
    linkStyle 2 stroke:#4682b4,stroke-width:2px,color:#4682b4
    linkStyle 3 stroke:#8a2be2,stroke-width:2px,color:#8a2be2
    linkStyle 4 stroke:#8a2be2,stroke-width:2px,color:#8a2be2
    linkStyle 5 stroke:#228b22,stroke-width:2px,color:#228b22
```

Таким образом, оркестратор создаёт агента, который в свою очередь имеет определённую вычислительную мощность - количество параллельно вычисляющихся задач, на каждую из которых выделяется исполнитель, решающий эту задачу. Worker - не отдельный класс или структура, а лишь метод агента.

> [!NOTE]
> По умолчанию COMPUTING_POWER = 1, это счётчик горутин. Worker - горутина, выполняющая задачу.

### Требования по ТЗ

**Оркестратор**:

1. Добавление вычисления арифметического выражения
```zsh
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
 "expression": <строка с выражение>
}'
```
Тело ответа:
```json
{
 "id": <уникальный идентификатор выражения>
}
```
Коды ответа:
- 201 - выражение принято для вычисления,
- 422 - невалидные данные,
- 500 - что-то пошло не так.

2. Получение списка выражений
```zsh
curl --location 'localhost/api/v1/expressions'
```
Тело ответа:
```json
{
    "expressions": [
        {
            "id": <идентификатор выражения>,
            "status": <статус вычисления выражения>,
            "result": <результат выражения>
        },
        {
            "id": <идентификатор выражения>,
            "status": <статус вычисления выражения>,
            "result": <результат выражения>
        }
    ]
}
```
Коды ответа:
- 200 - успешно получен список выражений,
- 500 - что-то пошло не так.

3. Получение выражения по его идентификатор
```zsh
curl --location 'localhost/api/v1/expressions/:id'
```
Тело ответа:
```json
{
    "expression":
        {
            "id": <идентификатор выражения>,
            "status": <статус вычисления выражения>,
            "result": <результат выражения>
        }
}
```
Коды ответа:
- 200 - успешно получено выражение,
- 404 - нет такого выражения,
- 500 - что-то пошло не так.

4. Получение задачи для выполнения
```zsh
curl --location 'localhost/internal/task'
```
Тело ответа
```json
{
    "task":
        {
            "id": <идентификатор задачи>,
            "arg1": <имя первого аргумента>,
            "arg2": <имя второго аргумента>,
            "operation": <операция>,
            "operation_time": <время выполнения операции>
        }
}
```
Коды ответа:
- 200 - успешно получена задача,
- 404 - нет задачи,
- 500 - что-то пошло не так.

5. Прием результата обработки данных.
```zsh
curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
  "id": 1,
  "result": 2.5
}'
```
Коды ответа:
- 200 - успешно записан результат,
- 404 - нет такой задачи,
- 422 - невалидные данные,
- 500 - что-то пошло не так.
