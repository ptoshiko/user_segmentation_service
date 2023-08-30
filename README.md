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