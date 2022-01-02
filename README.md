## Simple REST API with Go
## NO ORM

# Under Construction



### Tech Stack
 - Go (GoLang)
 - Gorilla Toolkit (mainly mux)
 - MySQL (storage)
 - JWT (for "stateless" authorization)
 - Redis** (jwt blacklisting)

### Development
 - Docker (Ubuntu 20.04)
 - Git (source control)
 - Github (of course ^_^ for remote code repo)
 - Visual Studio Code (Development Environment)


### Features
 - [x] Signup
 - [x] Login
 - [x] Logout
 - Task
	- Create
	- Read
	- Update
	- Delete**

### API Endpoints (So far)
| Feature | Endpoint | HTTP Verb | Description |
| ------- | -------- | --------- | ----------- |
| Signup | /auth/signup | POST | Create a user in the server that can log in. See below for required fields |
| Login | /auth/login | POST | Lets a registered user to authenticate. See below for required fields |
| Logout | /auth/logout | GET | Unauthenticate currently logged in user |


### Request Forms Fields And Constraints
Signup
```
{
  "firstname": <string: 1-16 characters: utf8>,
  "lastName": <string: 1-16 characters: utf8>,
  "email address": <string: 128 max characters: utf8: must conform to an email address>,
  "password": <string: 8 min characters: utf8>
}
```

Login
```
{
  "email address": <string: 128 max characters: utf8: must conform to an email address>,
  "password": <string: 8 min characters: utf8>
}
```