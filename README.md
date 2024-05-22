# Тестовое задание для стажёра Backend
## Cтэк
- **Язык:** Golang.
- **Базы данных:**  PostgreSQL, Redis.
- **Деплой зависимостей и  сервиса**: Docker и Docker Compose.

## Запуск
### Приложение может не запуститься с первого раза
```
make start 
```
## Тесты

```
make test #local
```
```
make test-docker #docker
```
## Интерфейс API реализован через swagger строго по данному в задании конфигу
### Swagger url
```
http://localhost:8000/swagger/index.html#/
```

## Дополнительные задания:
1. Так как фича и тэг по условию  - уникальное значение, было решено сделать индекс (feature_id,tag_id). Из-за этого скорость поиска повышается, однако  снижается скорость изменения. Также реализовано кэширование   Redis для метода  /banner_user   с временем жизни кэша 5 минут. С флагом use_last_revision пользователь получает некэшированные данные напрямую из PostgresQL.  
2. Нагрузочное тестирование не проведено
3. История изменений реализована. Так как по условию нужно смотреть последние 3 значения, у меня возникла дилемма, как  это реализовать. Я решил создать отдельную таблицу истории, а также метод, который с заданным интервалом оставляет только  3 последних изменения. Остальные записи удаляются.
4. Метод  удаления баннеров по фиче или тегу реализован. 
5. Реализовано mock-тестирование.
6. Линтер описан в [.golangci.yml](.golangci.yml).





# Примеры запросов
## Ниже приведены дублированние запросы из Swagger.
## Токен admin или user можно  получить по юрлу /login   с телом {"username":"admin"} или {"username":"notadmin"} соотвественно. Токен лежит в хедере в Authorization в виде "Bearer  eto_access_token". Чтобы делать запросы в Swagger, нужно добавить его в Authorize
![docs/img.png](docs/img.png)

## 1.Получение баннера для пользователя
### tag_id и feature_id обязательны. Со флагом use_last_revision пользователь получает некэшированные данные напрямую из PostgresQL.
### Запрос
```
curl -X 'GET' \
'http://localhost:8000/user_banner?tag_id=2&feature_id=3' \
-H 'accept: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2NzQ4MDYsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.3PJ9D8PZqhSMuRs8QzYoDG4ZvloGF_SE95_zdKTZVWg'
```
### Ответ
```
{
  "content": "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}"
}

```

## 2.Получение всех баннеров c фильтрацией по фиче и/или тегу
### Параметры tag_id, feature_id, limit, offset опциональны. Если ничего не указаать, то будут выведены все баннеры.
### Запрос
```
curl -X 'GET' \
'http://localhost:8000/banner?tag_id=2&feature_id=3' \
-H 'accept: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2NzQ4MDYsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.3PJ9D8PZqhSMuRs8QzYoDG4ZvloGF_SE95_zdKTZVWg'
```
### Ответ
```
[
  {
    "banner_id": 52,
    "tag_ids": [
      1,
      2
    ],
    "feature_id": 3,
    "content": "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
    "is_active": true,
    "created_at": "2024-04-12T23:17:17.825199Z",
    "updated_at": "2024-04-14T10:40:21.446315Z"
  }
]

```

## 3.Создание нового баннера
###  Все поля обязательны 
### Запрос
```
curl -X 'POST' \
  'http://localhost:8000/banner' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2OTIxMzAsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.PwJ-nnGI5TgrGYMhSzW2xJkxXnjXfpT022xefTpINEI' \
  -H 'Content-Type: application/json' \
  -d '{
      "content": "{}",
      "feature_id":4,
      "is_active": true,
      "tag_ids": [4]
}'
```
### Ответ
```
{
    "id":1
}

```

## 4.Обновление содержимого баннера
### Все поля опциональны. Также я решил, что tag_ids может быть нулевым(только пользователь не сможет найти такой баннер)
### Запрос
```
curl -X 'PATCH' \
'http://localhost:8000/banner/53' \
-H 'accept: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2ODM1MjAsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.3sqnpBYIJJjUIULrpG59jkKZ2kp_jVY8r4TRU6mcZ_I' \
-H 'Content-Type: application/json' \
-d '{
        "content": "{}",
        "feature_id": 4,
        "is_active": true,
        "tag_ids": []
}'

```
### Ответ
```
status code 201
```

## 5.Удаление баннера по идентификатору

### id обязателен. История баннера стирается.
### Запрос
```
curl -X 'DELETE' \
  'http://localhost:8000/banner/53' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2ODM1MjAsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.3sqnpBYIJJjUIULrpG59jkKZ2kp_jVY8r4TRU6mcZ_I'

```
### Ответ
```
status code 204
```
## 6.Удаление баннера по тэгу и фиче

### Поля  tag_id и feature_id обязательны. Возвращает ID удаленного банера. История баннера стирается.
### Запрос
```
curl -X 'DELETE' \
  'http://localhost:8000/banner?tag_id=2&feature_id=3' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2ODM1MjAsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.3sqnpBYIJJjUIULrpG59jkKZ2kp_jVY8r4TRU6mcZ_I'
```
### Ответ
```
{
    "id":1
}
```

## 7.История изменений баннера

### Параметр limit опционален. Возвращает не больше 3 последних изменений. Текущую конфигурацию не возвращает.
### Запрос
```
curl -X 'GET' \
  'http://localhost:8000/banner_history/56?limit=3' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2ODM1MjAsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.3sqnpBYIJJjUIULrpG59jkKZ2kp_jVY8r4TRU6mcZ_I'
```
### Ответ
```
[
  {
    "banner_id": 56,
    "tag_ids": [
      3
    ],
    "feature_id": 1,
    "content": "{\"title\": \"some_title\"}",
    "is_active": true,
    "created_at": "2024-04-14T13:44:41.473662Z",
    "updated_at": "2024-04-14T14:08:40.831097Z"
  },
  {
    "banner_id": 56,
    "tag_ids": [
      3
    ],
    "feature_id": 1,
    "content": "{\"slatt\": \"slatt\"}",
    "is_active": true,
    "created_at": "2024-04-14T13:44:41.473662Z",
    "updated_at": "2024-04-14T14:08:09.860374Z"
  },
  {
    "banner_id": 56,
    "tag_ids": [],
    "feature_id": 1,
    "content": "{\"slatt\": \"slatt\"}",
    "is_active": true,
    "created_at": "2024-04-14T13:44:41.473662Z",
    "updated_at": "2024-04-14T13:57:51.873896Z"
  }
]
```

## 8.Получение токена

### Предосталвяет соответствующие права admin или notadmin  при указании их в username.
### Запрос
```
curl -X 'POST' \
  'http://localhost:8000/login' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2ODM1MjAsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.3sqnpBYIJJjUIULrpG59jkKZ2kp_jVY8r4TRU6mcZ_I' \
  -H 'Content-Type: application/json' \
  -d    '{
          "username": "admin"
        }'
```
### Ответ
```
{
  "access": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU2ODUxOTIsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIiwiYWRtaW4iOnRydWV9.d-RfnChiqzdPlrTB1zryq5wnI0NnMXhx4BNk1YlPfJo"
}
```

## Структура таблиц

```
create table tag(
    id bigint primary key
);
create table feature(
    id bigint primary key
);
create table banner(
    id serial primary key ,
    content text,
    is_active boolean,
    created_at timestamp default  now() not null,
    updated_at timestamp default  now()  not null
);
create table banner_tag(
    banner_id  integer references banner(id) on DELETE cascade ,
    tag_id  bigint references tag(id) on DELETE cascade ,
    feature_id  bigint references feature(id) on DELETE cascade ,
    constraint banner_tag_feature_uc unique (tag_id,feature_id)
);
create table banner_history(
    id serial primary key ,
    banner_id integer references banner(id) on DELETE cascade,
    content text,
    is_active boolean,
    created_at timestamp  not null,
    updated_at timestamp  not null
);

create table banner_tag_history(
    banner_history_id  integer references banner_history(id) on delete cascade ,
    tag_id  bigint references tag(id) ,
    feature_id bigint references feature(id),
    constraint banner_tag_history_feature_uc unique (banner_history_id,tag_id,feature_id)
);


```
