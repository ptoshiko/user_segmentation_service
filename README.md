# Тестовое задание для Тестовое задание для стажёра Backend Avito
# Сервис динамического сегментирования пользователей

### Описание задачи

Требуется реализовать сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент). 

Полное описание задачи по [ссылке](https://github.com/avito-tech/backend-trainee-assignment-2023).

### Технические требования 

1. Сервис должен предоставлять HTTP API с форматом JSON как при отправке запроса, так и при получении результата.
2. Язык разработки: Golang.
3. Фреймворки и библиотеки можно использовать любые.
4. Реляционная СУБД: MySQL или PostgreSQL.
5. Использование docker и docker-compose для поднятия и развертывания dev-среды.

### Запуск
```bash
$ git clone https://github.com/ptoshiko/avito_assignment.git
$ cd avito_assignment
$ make run
```

### Запуск тестов 
```bash
make testdb
nake t
```

### Реализованный функционал 

#### CreateSegment  
Метод создания сегмента.Принимает slug (название) сегмента.

**Пример запроса** 

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "seg_name": "AVITO_DISCOUNT_50"           
}' http://localhost:8080/segment/create

```
<br>

**Ответ**
```json
{
  "message":"Segment created successfully"
}
```

#### DeleteSegment 
Метод удаления сегмента. Принимает slug (название) сегмента.

**Пример запроса** 

```bash
curl -X DELETE -H "Content-Type: application/json" -d '{
  "seg_name": "AVITO_DISCOUNT_30"
}' http://localhost:8080/segment/delete

```
<br>

**Ответ**
```json
{
  "message": "Segment deleted successfully"
}
```

#### UpdateUserSegments 
Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.


Валидным считается запрос, в котором:
1. сегменты из списка для добавления присутвуют в таблице сегментов 
2. сегменты из списка для удаления принадлежат пользователю
3. один из списков может быть пустым 
В остальных случаях запрос считается невалидным. 

**Пример запроса** 

```bash
curl -X PATCH -H "Content-Type: application/json" -d '{
  "user_id": 1,
  "segments_to_add": ["AVITO_VOICE_MESSAGES"],
  "segments_to_remove": ["AVITO_PERFORMANCE_VAS"]
}' http://localhost:8080/user/1

```
<br>

**Ответ**
```json
{
  "message": "User segments updated successfully"
}
```

#### GetUserSegments
Метод получения активных сегментов пользователя. Принимает на вход id пользователя.

```bash
curl -X GET http://localhost:8080/user/1

```
<br>

**Ответ**
```json
{
  [{"seg_id":1,"seg_name":"AVITO_DISCOUNT_30"},{"seg_id":3,"seg_name":"AVITO_VOICE_MESSAGES"}]
}
```




