# Тестовое задание Junior Backend Developer

**Используемые технологии:**

- Go
- JWT
- MongoDB

**Задание:**

Написать часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя сидентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refreshтокенов

**Требования:**

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базеисключительно в виде bcrypt хеша, должен быть защищен от изменения настороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

**Результат:**

Результат выполнения задания нужно предоставить в виде исходного кода на Github.

# О решении

## Как запустить
 
1. Собрать Docker контейнер сервера:
```bash
docker build -t testjunior:v1 .
```
2. При необходимости изменить переменные окружения в docker-compose.yml
3. Запустить Docker Compose:
```bash
docker-compose up
```

## Как использовать

- Получение токенов
```
/token/generate POST
```
Пример тела запроса:
```json
{
    "user_guid": "f022c263-4796-4bdc-92ac-4c74020183b3"
}
```
Пример тела ответа:
```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMyMDE3NTAsInVzZXJfZ3VpZCI6ImYwMjJjMjYzLTQ3OTYtNGJkYy05MmFjLTRjNzQwMjAxODNiMyJ9.MBmJ2HX5VM0XgEEqwk7eWpteOHgKOzNoWBhL2zyHTu3DkPuPmTujCjYcnC5elCdNpKyMyWasOiXO67kZbgk67Q",
    "refresh_token": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk9UTXlPRGM0TlRBc0ltbGtJam9pWmpRd1pqWmpNek10TUdabE1DMDBNVGxsTFRnNVpETXRaR0ZoTnpSbVlXWTBZalV4SWl3aWRYTmxjbDluZFdsa0lqb2laakF5TW1NeU5qTXRORGM1TmkwMFltUmpMVGt5WVdNdE5HTTNOREF5TURFNE0ySXpJbjAuUEh3Z051d3BUeUwteGU1SGJPRXA1VGF3QnIySFpvdGtOajMzNDB0OXdfNXFhTF9zcVo2OVVrU182UzdXN3RMSlVtZjlVTVBGTHVkcTl0ZEktRnZYUmc="
}
```
- Обновление токена
```
/token/refresh POST
```
Пример тела запроса:
```json
{
    "refresh_token": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk9UTXlPRGM0TlRBc0ltbGtJam9pWmpRd1pqWmpNek10TUdabE1DMDBNVGxsTFRnNVpETXRaR0ZoTnpSbVlXWTBZalV4SWl3aWRYTmxjbDluZFdsa0lqb2laakF5TW1NeU5qTXRORGM1TmkwMFltUmpMVGt5WVdNdE5HTTNOREF5TURFNE0ySXpJbjAuUEh3Z051d3BUeUwteGU1SGJPRXA1VGF3QnIySFpvdGtOajMzNDB0OXdfNXFhTF9zcVo2OVVrU182UzdXN3RMSlVtZjlVTVBGTHVkcTl0ZEktRnZYUmc="
}
```
Пример тела ответа:
```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMyMDI0NDYsInVzZXJfZ3VpZCI6ImYwMjJjMjYzLTQ3OTYtNGJkYy05MmFjLTRjNzQwMjAxODNiMyJ9.h5iXOFNxuGg29n0PWkhURbEfuYLgEw9qhKnAuCNAz_PYtQ1iwk1AESXVLcjkQrKxhknTesxSPzaIYP174I1ijw",
    "refresh_token": "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk9UTXlPRGcxTkRZc0ltbGtJam9pWWpkaU16STRNR0V0WldFM015MDBZMk5oTFdJMllqSXRNVFZpTldKbE5qSmtabVJqSWl3aWRYTmxjbDluZFdsa0lqb2laakF5TW1NeU5qTXRORGM1TmkwMFltUmpMVGt5WVdNdE5HTTNOREF5TURFNE0ySXpJbjAubUxJTW9XZXpvTGRPaUkyMEY3eVN4dXlSNGJ3RUZmSXNJN0VMSHhleUc2ZFFFYmU2bGozaFM3NVhLU3pCNEU1YTFKSDVKR3BydGFiMWpZdms3T1Z3dGc="
}
```
## О переменных окружения
- SERVER_PORT - порт, который будет слушать сервер
- DB_URI - URI для подключения к MongoDB в формате mongodb://user:password@host:port/
- LOG_FILE - файл для записи логов
- DB_NAME - название базы данных, в которой находится (или будет находиться) коллекция с данными о refresh-токенах
- ACCESS_TOKEN_TTL - время жизни access токена в формате 1234{единица измерения} (единицы измерения ns - наносекунды, ms - миллисекунды, s - секунды, m - минуты, h - часы )
- REFRESH_TOKEN_TTL: время жизни refresh токена (формат такой же)
- ACCESS_TOKEN_SECRET: секретный ключ, которым будет подписываться JWT access токена,
- REFRESH_TOKEN_SECRET: секретный ключ, которым будет подписываться JWT refresh токена
- TEST_DB_URI - URI базы данных для запуска тестов (go test)