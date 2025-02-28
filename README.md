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

### Как это работает?

```graph TD
U[Пользователь] -->|POST /calculate| O[Оркестратор]
U -->|GET /expressions| O
O -->|GET /internal/task| A[Агент]
subgraph Workers
W1[Worker 1]
W2[Worker 2]
W3[Worker 3]
Wn[...]
end
A -->|Computing power| Workers
Workers -->|POST /internal/task| A
A -->|POST /internal/task| O
```