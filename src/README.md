# Productivity App Backend

## This project provides a number of APIs for use in goal-tracking. The project architecture is microservice.

1. User API

*Create User*: Create a new user for the application
 
 `POST /users`
```json
{ 
"first_name":"string",
"last_name":"string"
}
```
Response:

Status `201`
```json
{
  "first_name": "string",
  "last_name": "string", 
  "id": "object_id"
}
```

*Read Single User by Id*: Retrieve the details of a user with id `id`

`GET /users/{id}`

Response

Status `200`

```json
{
"first_name": "string",
"last_name": "string",
"id": "object_id"
}
``` 

 
*Read All Users* Retrieve a list of all users registered in the application

`GET /users`

Response

Status `200`

```json
[
  {
    "first_name": "string",
    "last_name": "string",
    "id": "object_id"
  }
]
```

*Update User* Update specified details of a user by passing in the user's `id` and the key-value pairs to change

`PATCH /users/{id}`

```json
"first_name": "string" (optional),
"last_name": "string" (optional)
```

Response

Status `200`
```json
"first_name": "string",
"last_name": "string",
"id": "object_id"
```

*Delete User* Remove user from application by specifying their user `id`.

`DELETE /users/{id}`

Response

Status `204`

```json
no content
```

2. Goal API

*Create Goal*

3. Event API

4. Data Summary API