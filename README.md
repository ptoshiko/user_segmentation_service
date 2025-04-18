# Test Assignment for Avito Backend Internship
# User Dynamic Segmentation Service

[Версия на русском](README_RU.md)

### Task Description

The task is to implement a service that stores users and the segments they belong to (creating, modifying, deleting segments, as well as adding and removing users from segments).

Full task description is available at this [link](https://github.com/avito-tech/backend-trainee-assignment-2023).

### Technical Requirements

1. The service must provide an HTTP API with JSON format for both sending requests and receiving results.
2. Development language: Golang.
3. Any frameworks and libraries can be used.
4. Relational DBMS: MySQL or PostgreSQL.
5. Use of docker and docker-compose for setting up and deploying the dev environment.

### Setup
```bash
$ git clone https://github.com/ptoshiko/avito_assignment.git
$ cd avito_assignment
$ make run
```

### Running Tests
```bash
make testdb
make t
```

### Implemented Functionality

#### CreateSegment
Method for creating a segment. Accepts a segment slug (name).

**Request Example**

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "seg_name": "AVITO_DISCOUNT_50"           
}' http://localhost:8080/segment/create

```
<br>

**Response**
```json
{
  "message":"Segment created successfully"
}
```

#### DeleteSegment
Method for deleting a segment. Accepts a segment slug (name).

**Request Example**

```bash
curl -X DELETE -H "Content-Type: application/json" -d '{
  "seg_name": "AVITO_DISCOUNT_30"
}' http://localhost:8080/segment/delete

```
<br>

**Response**
```json
{
  "message": "Segment deleted successfully"
}
```

#### UpdateUserSegments
Method for adding a user to a segment. Accepts a list of segment slugs (names) to add to the user, a list of segment slugs (names) to remove from the user, and a user ID.

A request is considered valid if:
1. Segments from the addition list are present in the segments table
2. Segments from the removal list belong to the user
3. One of the lists can be empty
In other cases, the request is considered invalid.

**Request Example**

```bash
curl -X PATCH -H "Content-Type: application/json" -d '{
  "user_id": 1,
  "segments_to_add": ["AVITO_VOICE_MESSAGES"],
  "segments_to_remove": ["AVITO_PERFORMANCE_VAS"]
}' http://localhost:8080/user/1

```
<br>

**Response**
```json
{
  "message": "User segments updated successfully"
}
```

#### GetUserSegments
Method for retrieving active segments of a user. Takes a user ID as input.

```bash
curl -X GET http://localhost:8080/user/1

```
<br>

**Response**
```json
{
  [{"seg_id":1,"seg_name":"AVITO_DISCOUNT_30"},{"seg_id":3,"seg_name":"AVITO_VOICE_MESSAGES"}]
}
```



