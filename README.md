### Примеры запросов

curl -X DELETE -H "Content-Type: application/json" -d '{
  "seg_name": "Spring"
}' http://localhost:8080/segment/1

curl -X POST -H "Content-Type: application/json" -d '{
  "seg_name": "Samuel"           
}' http://localhost:8080/segment/create

curl -X PATCH -H "Content-Type: application/json" -d '{
  "user_id": 1,
  "segments_to_add": ["Azat"],
  "segments_to_remove": []
}' http://localhost:8080/user/1

curl -X GET http://localhost:8080/user/1

расписать про update 

Метод обновления. валидным считается запрос в котором передаются сегменты которые:
1. в случае добавления сегментов сегменты присутвуют в таблице сегментов 
2. в случае удаление сегментов сегменты принадлежат пользователю
3. один из списков может быть пустым 
В остальных случаях запрос считается невалидным