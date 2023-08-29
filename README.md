### Примеры запросов

curl -X POST -H "Content-Type: application/json" -d '{
  "seg_name": "AVITO_DISCOUNT_50"
}' http://127.0.0.1:8080/create_segment

curl -X POST -H "Content-Type: application/json" -d '{
  "seg_name": "AVITO_PERFORMANCE_VAS"
}' http://127.0.0.1:8080/delete_segment

curl -X POST -H "Content-Type: application/json" -d '{
  "id": 1
}' http://127.0.0.1:8080/get_user_segments


curl -X POST -H "Content-Type: application/json" -d '{
  "user_id": "1",
  "segments_to_add": ["AVITO_PERFORMANCE_VAS"],         
  "segments_to_remove": ["AVITO_DISCOUNT_30"]   
}' http://localhost:8080/update_user_segments


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