# follow-service

## Overview
This service handles all user-follow relationships and is part of the two-service set up for **Backend Intsernship at Pratipili**

It exposes endpoints to:
- Create and list users
- Follow and unfollow users
- List followers and followings of a given user

The **GraphQL Service (gateway)** consumes these APIs to expose GraphQL queries and mutations.

## For a detailed tech stack choices, deployment steps, etc, refer graphql-service repo
Link: https://github.com/Infamous003/graphql-service

## Base URL
Link: https://follow-service-n9fk.onrender.com

## API Documentation

#### User Routes

`POST /users`
Creates a new user
```json
{"username": "testuser123"}
```

```json
{
	"user": {
		"ID": 5,
		"Username": "testuser123",
		"CreatedAt": "2025-12-10T17:43:32.136748Z"
	}
}
```

`Get /users`
List all users

```json
{
	"users": [
		{
			"ID": 1,
			"Username": "infamous03",
			"CreatedAt": "2025-12-10T11:05:26.499475Z"
		},
		{
			"ID": 2,
			"Username": "syedmehdi03",
			"CreatedAt": "2025-12-10T11:08:40.128579Z"
		}
	]
}
```

#### Follow Routes
`POST /follow`

Request Body:
```json
{
    "follower_id": 1,
    "followee_id": 2
}
```

Response (success)
```json
{
    "message": "followed successfully"
}
```
Response (already following)
```json
{
    "message": "already following this user"
}
```
Response (self following)
```json
{
    "message": "cannot follow yourself"
}
```

`POST unfollow`

Request Body
```json
{
    "follower_id": 1,
    "followee_id": 2
}
```

Response
```json
{ 
    "message": "unfollowed successfully"
}
```

Response:
```json
{ 
    "message": "the requested resource could not be found"
}
```

`GET /users/{id}/followers`
Get all followers of a specific user
```json
{
	"followers": [
		{
			"ID": 1,
			"Username": "infamous03",
			"CreatedAt": "2025-12-10T11:05:26.499475Z"
		},
		{
			"ID": 2,
			"Username": "syedmehdi03",
			"CreatedAt": "2025-12-10T11:08:40.128579Z"
		}
	]
}
```

`GET /users/{id}/following`
Get all followings of a user
```json
{
	"following": [
		{
			"ID": 3,
			"Username": "arfathsyed",
			"CreatedAt": "2025-12-10T11:08:48.024457Z"
		},
		{
			"ID": 4,
			"Username": "syed03",
			"CreatedAt": "2025-12-10T11:09:03.090494Z"
		}
	]
}
```

## Run Locally
```bash
git clone https://github.com/Infamous003/follow-service.git
cd follow-service
make db/migrations/up
make run
```
Make sure Postgres is up and running and env variables are set

## Deployment
Deployed on **Render.com**
URL: https://follow-service-n9fk.onrender.com

