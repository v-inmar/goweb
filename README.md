## Simple REST API with Go and Gorilla Mux
## With MySQL

### NO ORM

**Under Construction**

Just a simple rest API for a Todo List


### API Endpoints

| URI   |      Verb      |  Description |
|----------|:-------------:|------:|
| /todos |  POST | Creates a new todo item |
| /todos |  GET | Retrieve all the todo items |
| /todos/{id} |  GET | Retrieve a single todo item |
| /todos/{id} |  PUT | Update a todo item |
| /todos/{id} |  DEL | Delete a todo item |


#### Create new todo item request body signature

```json
{
	"title": "<string value: cannot empty or missing>",
	"body": "<string value>: can be empty"
}

```

#### Update todo item request body signature

```json
{
	"title": "<string value: will be ignored if empty or missing>",
	"body": "<string value>: can be empty"
}

```


### Bugs :poop: (On-going)

| ID   |      API      |  Description |
|----------|:-------------:|------:|
| 1 |  /todos/{id}:PUT | Missing json key 'body' will not ignore but remove current body value |
