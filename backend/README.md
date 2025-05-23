## Установка
Устанавливаются с помощью docker compose следующей командой:
```bash
docker compose up --build
```
Для последующих запусков без пересборки контейнеров:
```bash
docker compose up
```
Для того, чтобы запустить контейнеры на заднем фоне:
```bash
docker compose up -d
```
Но тогда в окне терминала не будут выводиться логи, для их вывода можно использовать
```bash
docker compose logs --tail=10 # выведет последние 10 строк логов
```
Соберутся и запустятся 2 контейнера: серверное приложение и база данных postgres.
Для остановки контейнеров:
```bash
docker compose down
```
Изменения в базе данных будут сохраняться между перезапусками и пересборками контейнеров.
Чтобы вернуть базу данных к начальному состоянию:
```bash
docker compose down -v
```
Это удалит volume - хранилище данных для БД. При следующем запуске docker compose БД будет заново проинициализирована начальными значениями.
## Использование
Для тестирования в VS code с расширением "REST Сlient" можно использовать файл ```rest_api.http```, там есть уже готовые запросы в нужном формате.

При добавлении, изменении, удалении всех ресурсов будет возвращен json созданного, обновленного, удаленного ресурса.

### Статусы ответов:

200 (ок) - при успешном выполнении запроса,

201 (created) - при создании ресурса,

400 (bad request) - неверный формат запроса,

404 (not found) - при отсутствии ресурса,

403 (forbidden) - при недостаточных правах на доступ к ресурсам,

401 (unauthorized) - при отсутствии, истечении или невалидном токене, а также при неверном пароле,

409 (conflict) - нарушение уникальности (например создание нескольких отзывов одного и того же пользователя на один и тот же смартфон, или добавление в корзину смартфона, который там уже есть),

500 (internal server error) - внутренняя ошибка сервера.

Во всех запросах request header Content-Type можно не указывать, приложение все равно будет относиться к содержимому как к json

Для защищенных эндпойнтов должен быть установлен header Authorization куда помещается jwt token, возвращаемый при логине. Допускается формат поля как ```Bearer {token}```, так и просто ```{token}```. Доступ имеют либо владелец ресурса, либо пользователь с ролью ```admin```.

Токен действителен в течении 24 часов.
## Запросы
### Получить все смартфоны:
```
GET "http://localhost:8081/api/v1/smartphones"
```

### Для получения лишь смартфонов с определенными айди:
```
GET "http://localhost:8081/api/v1/smartphones?ids=1,3,4"
```
### Получить один смартфон с определенным айди:
```
GET "http://localhost:8081/api/v1/smartphones/{smartphone_id}"
```
Пример:
```
GET "http://localhost:8081/api/v1/smartphones/1"
```
### Поля смарфтона:
```json
{
  "id": 1,
  "model": "iPhone 16",
  "producer": "Apple",
  "memory": 128,
  "ram": 8,
  "display_size": 6.1,
  "price": 999,
  "ratings_sum": 9,
  "ratings_count": 2,
  "image_path": "https://c.dns-shop.ru/thumb/st1/fit/0/0/1043f341d851923dda2ac92e50f089a1/14ce8c6a5fbaef30feb3cb6b7d742546c045c44eb9207be4acec68cade72a7cf.jpg.webp",
  "description": "Introducing the all-new iPhone 16 where innovation meets elegance. With a sleek design and cutting-edge technology, the iPhone 16 delivers a stunning display, incredible camera capabilities, and lightning-fast performance that transforms the way you experience mobile devices. Whether you are capturing memories in breathtaking detail, enjoying your favorite content in vibrant color, or seamlessly multitasking, the iPhone 16 elevates every moment. Powered by Apple most advanced chipset, it brings unmatched speed, efficiency, and security to your fingertips. Step into the future with the iPhone 16 - a new standard in smartphone excellence.",
  "reviews": [
    {
      "id": 1,
      "smartphone_id": 1,
      "user_id": 2,
      "user_name": "user1",
      "rating": 5,
      "comment": "Best as always.",
      "created_at": "2025-05-21T19:50:51.888096Z",
      "updated_at": "2025-05-21T19:50:51.888096Z"
    },
    {
      "id": 2,
      "smartphone_id": 1,
      "user_id": 3,
      "user_name": "user2",
      "rating": 4,
      "created_at": "2025-05-21T19:50:51.888096Z",
      "updated_at": "2025-05-21T19:50:51.888096Z"
    }
  ]
}
```
В запросах нескольких смартфонов ```api/v1/smartphones``` поле ```reviews``` будет полностью отсутствовать
### Регистрации нового пользователя:
```
POST "http://localhost:8081/api/v1/signup"
Content-Type: application/json

{
    "name": "username",
    "password": "password"
}
```
### Логин:
```
POST "http://localhost:8081/api/v1/login"
Content-Type: application/json

{
    "name": "username",
    "password": "password"
}
```
### Возвращаемое значение при успешном логине:
```json
{
  "user": {
    "id": 1,
    "name": "admin",
    "role": "admin",
    "created_at": "2025-03-21T15:10:01.619971Z",
    "cart": {
      "id": 7,
      "user_id": 1,
      "created_at": "2025-03-21T15:10:01.619971Z",
      "updated_at": "2025-03-21T15:10:01.619971Z",
      "items": [
        {
          "id": 4,
          "cart_id": 7,
          "smartphone_id": 1,
          "quantity": 1
        },
        {
          "id": 7,
          "cart_id": 7,
          "smartphone_id": 2,
          "quantity": 1
        }
      ]
    }
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk"
}
```
```user``` содержит основную информацию о пользователе. ```cart``` - его корзина. Для всех новых пользователей автоматически будет создана корзина. ```items``` - содержимое корзины. ```token``` - jwt токен.

Токен можно декодировать, в нем содержатся поля ```sub``` - user_id, ```iss``` - "Smartbuy", ```iat``` - дата выпуска токена, ```exp``` - дата истечения токена, ```role```- роль (```admin``` или ```user```).

Роли пользователей: ```admin``` и ```user```. Пользователь с ролью ```admin``` имеет доступ ко всем защищенным эндпойнтам. В БД уже существует пользователь с именем ```admin``` и паролем ```admin``` для тестирования. У остальных пользователей пароль идентичен имени. Пароли хранятся в хешированном виде (bcrypt, 10)

### Получить всех пользователей:
```
GET http://localhost:8081/api/v1/users/
Authorization: {token}
```
Только для админов, поле ```cart``` будет отсутствовать
### Получить пользователя по айди:
```
GET http://localhost:8081/api/v1/users/{user_id}
Authorization: {token}
```
### Получить корзину по айди корзины:
```
GET http://localhost:8081/api/v1/carts/{cart_id}
Authorization: {token}
```
Пример:
```
GET http://localhost:8081/api/v1/carts/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk
```
### Поля корзины:
```json
{
  "id": 10,
  "created_at": "2025-03-21T15:13:54.800657Z",
  "updated_at": "2025-03-21T15:13:54.800657Z",
  "items": [
    {
      "id": 3,
      "smartphone_id": 1,
      "quantity": 1
    }
  ]
}
```
### Получить корзину по айди пользователя:
```
GET http://localhost:8081/api/v1/carts?user_id={user_id}
Authorization: {token}
```
Пример:
```
GET http://localhost:8081/api/v1/carts?user_id=1
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk
```
### Получить все корзины:
```
GET http://localhost:8081/api/v1/carts
Authorization: {token}
```
Только для админов, поле ```items``` будет отсутствовать
### Получить предметы в корзине по айди корзины:
```
GET http://localhost:8081/api/v1/carts/{cart_id}/items
Authorization: {token}
```
Пример:
```
GET http://localhost:8081/api/v1/carts/4/items
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk
```
### Поля предметов в корзине:
```json
[
  {
    "id": 3,
    "smartphone_id": 1,
    "quantity": 1
  },
  {
    "id": 5,
    "smartphone_id": 2,
    "quantity": 3
  }
]
```
### Добавить предмет в корзину по айди корзины:
```
POST "http://localhost:8081/api/v1/carts/{cart_id}/items"
Authorization: {token}

{
    "smartphone_id": 1
}
```
```quantity``` будет равно единице
### Изменить количество предмета в корзине:
```
PATCH "http://localhost:8081/api/v1/carts/{cart_id}/items/{item_id}"
Authorization: {token}

{
    "quantity": 3
}
```
Пример:
```
PATCH "http://localhost:8081/api/v1/carts/1/items/1"
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk

{
    "quantity": 3
}
```
### Удалить предмет из корзины:
```
DELETE http://localhost:8081/api/v1/carts/{cart_id}/items/{item_id}
Authorization: {token}
```
### Получить отзывы к смартфону
```
GET http://localhost:8081/api/v1/smartphones/{smartphone_id}/reviews
```
Пример:
```
GET http://localhost:8081/api/v1/smartphones/1/reviews
```
### Поля отзывов к смартфону:
```json
[
  {
    "id": 1,
    "smartphone_id": 1,
    "user_id": 2,
    "user_name": "user1",
    "rating": 5,
    "comment": "Best as always.",
    "created_at": "2025-05-21T19:50:51.888096Z",
    "updated_at": "2025-05-21T19:50:51.888096Z"
  },
  {
    "id": 2,
    "smartphone_id": 1,
    "user_id": 3,
    "user_name": "user2",
    "rating": 4,
    "created_at": "2025-05-21T19:50:51.888096Z",
    "updated_at": "2025-05-21T19:50:51.888096Z"
  }
]
```
Поле ```comment``` может отсутствовать
Поле ```user_name``` присутствует только при GET всех отзывов по айди смартфона, в остальных отсутствует.
### Добавить отзыв к смартфону:
```
POST "http://localhost:8081/api/v1/smartphones/{smartphone_id}/reviews"
Authorization: {token}

{
    "rating": 4,
    "comment": "optional comment" // может отсутствовать
}
```
Примеры:
```
POST "http://localhost:8081/api/v1/smartphones/1/reviews"
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk

{
    "rating": 5,
    "comment": "great smartphone"
}
```
```
POST "http://localhost:8081/api/v1/smartphones/1/reviews"
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk

{
    "rating": 3
}
```
### Изменить отзыв:
```
PATCH "http://localhost:8081/api/v1/smartphones/{smartphone_id}/reviews/{review_id}"
Authorization: {token}

{
    "rating": 5,
    "comment": "new optional comment" // может отсутствовать
}
```
Примеры:
```
PATCH "http://localhost:8081/api/v1/smartphones/1/reviews/5"
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk

{
    "rating": 4,
    "comment": "decent phone"
}
```
```
PATCH "http://localhost:8081/api/v1/smartphones/1/reviews/5"
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJpc3MiOiJTbWFydGJ1eSIsInN1YiI6IjEiLCJleHAiOjE3NDI2Njk2NzUsImlhdCI6MTc0MjU4MzI3NX0.vUuOxfPBbsD30xg5bQdsNC4Dti2UTkPA06icmGvnKNk

{
    "rating": 3
}
```
### Удалить отзыв:
```
DELETE http://localhost:8081/api/v1/smartphones/{smartphone_id}/reviews/{review_id}
Authorization: {token}
```
