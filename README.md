Микросервис для работы с балансом пользователей

Для работы с данной реализацией необходимо установить библиотеки
"github.com/joho/godotenv"

"github.com/lib/pq"

"github.com/jmoiron/sqlx"

"github.com/sirupsen/logrus"

"github.com/spf13/viper"

"github.com/gin-gonic/gin"

Необходимо запустить докер контейнер  docker run --name=money-service -e POSTGRES_PASSWORD='qwerty' -p 5437:5432 -d --rm postgres

Далее подключиться к бд docker exec -it имя_контейнера /bin/bash

Следующий шаг применить утилитой мигрейт команду для создания файлов миграций 

migrate create -ext sql -dir ./schema -seq init

Дальше запускаем миграции
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5437/postgres?sslmode=disable' up


Для запуска приложения - нужно написать в консоли   go run cmd/main.go

1)router.GET("/getBalance", h.GetBalance) для получения баланса
2)router.GET("/addMoney", h.AddMoney) для добавления денег юзеру 
3)router.GET("/reserveMoney", h.ReserveMoney) резервация денег на отдельный счет
4)router.GET("/approveDeal", h.DealSuccess) разрезервация денег и отправка их в отчет
5)router.GET("/cancelOrder", h.CancelOrder) отмена резервации
6)router.GET("/getHistory", h.UserHistory) история баланса пользователя
7)router.GET("/getMonthReport", h.MonthReport) получения отчета по заданному году и месяцу

Примеры входных данных и ответов , использовался постман 
1
```yaml
Input:
{
    "id":2
}
Output:
{
    "balance": 3002
}
Input:
{
    "id":1
}
Output:
{
    "balance": 400.12
}


2
Input:
{
    "id":1,
    "money_amount":100.12
}
Output:
{
    "balance": 100.12
}


Input:
{
    "id":1,
    "money_amount":300
}
Output:
{
    "balance": 400.12
}

Input:
{
    "id":2,
    "money_amount":300
}
Output:
{
    "balance": 300
}


3
Input:
{
    "user_id":1,
    "service_id":1,
    "order_id":1,
    "order_cost":100.12
}
Output:
{
    "amount_reserved": 100.12,
    "new_reserved_balance": 100.12
}
Input:
{
    "user_id":1,
    "service_id":1,
    "order_id":2,
    "order_cost":10
}
Output:
{
    "amount_reserved": 10,
    "new_reserved_balance": 110.12
}
Input:
{
    "user_id":2,
    "service_id":1,
    "order_id":3,
    "order_cost":101.123
}
Output:
{
    "amount_reserved": 101.123,
    "new_reserved_balance": 101.123
}

4
Input:
{
    "id":1,
    "service_id":1,
    "order_id":1,
    "order_cost":100.12
}
Output:
{
    "deal_status": true
}
Input:
{
    "id":2,
    "service_id":1,
    "order_id":3,
    "order_cost":101.123
}
Output:
{
    "deal_status": true
}

5
Input:
{
    "id":1,
    "service_id":1,
    "order_id":2,
    "order_cost":10
}
Output:
{
    "deal_canceled": true
}
6
Input:
{
    "id":1
}
Output:
{
    "1": {
        "Event": "adding money to balance                 ",
        "Amount": 100.12,
        "Date": "2022-11-09T23:37:45.740039Z"
    },
    "2": {
        "Event": "adding money to balance                 ",
        "Amount": 300,
        "Date": "2022-11-09T23:37:52.318606Z"
    },
    "3": {
        "Event": "reserving money from balance            ",
        "Amount": 100.12,
        "Date": "2022-11-09T23:38:32.543114Z"
    },
    "4": {
        "Event": "reserving money from balance            ",
        "Amount": 10,
        "Date": "2022-11-09T23:38:40.041821Z"
    },
    "5": {
        "Event": "successful payment                      ",
        "Amount": 100.12,
        "Date": "2022-11-09T23:39:10.73423Z"
    },
    "6": {
        "Event": "canceling order                         ",
        "Amount": 10,
        "Date": "2022-11-09T23:41:25.597948Z"
    }
}
Input:
{
   "id":2
}
Output:
{
    "1": {
        "Event": "adding money to balance                 ",
        "Amount": 300,
        "Date": "2022-11-09T23:38:06.964693Z"
    },
    "2": {
        "Event": "reserving money from balance            ",
        "Amount": 101.123,
        "Date": "2022-11-09T23:38:50.941285Z"
    },
    "3": {
        "Event": "successful payment                      ",
        "Amount": 101.123,
        "Date": "2022-11-09T23:40:04.863109Z"
    }
}

7
Input:
{
    "date":"2022-11"
}
Output:
{
    "report_url": "./report.csv"
}
Возвращает ссылку на локальный файл репорт в котором лежит отчет
service_id: 1,total_money_amount: 201.243

