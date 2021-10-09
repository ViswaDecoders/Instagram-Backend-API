# Instagram-Backend-API(NO EXTRA DEPENDENCIES)

Designed and Developed an HTTP JSON REST API mocking of instagram posts and users schema which capable of the following structure operations given below. 
[*Golang*](https://golang.org/)(https://tour.golang.org/) is used for the API and [*MongoDB*](https://docs.mongodb.com/manual/tutorial/) is used as a storage. ***Documentation*** has been done where ever is necessary.

# Structure
- [x] Create an User
* Should be a POST request
* Use JSON request body
* URL should be ‘/users'
* e.g ``` curl.exe http://localhost:10000/users -H "Content-Type:application/json" -X POST -d '{\"id\":45,\"name\":\"Ramu\",\"email\":\"ram1@yahoo.com\",\"password\":\"rAmU4521@#\"}' ``` 
* Note: Pass will be encrypted at the time of storage

- [x] Get a user using id
* Should be a GET request
* Id should be in the url parameter
* URL should be ‘/users/{id here}’
* e.g ``` curl.exe http://localhost:10000/users/<id> ```
  
- [x] Create a Post
* Should be a POST request
* Use JSON request body
* URL should be ‘/posts'
* e.g ``` curl.exe http://localhost:10000/posts -H "Content-Type:application/json" -X POST -d '{\"id\":21,\"userif\":45,\"imageurl\":\"http://c:/documents/j.jpg\",\"caption\":\"scenary\"}' ``` 
* Note: Time stamp get atomatically taken from system and stored in MongoDB storage

  
- [x] Get a post using id
* Should be a GET request
* Id should be in the url parameter
* URL should be ‘/posts/{id here}
* e.g ``` curl.exe http://localhost:10000/posts/<id> ```
  
- [x] List all posts of a perticular user
* Should be a GET request
* URL should be ‘/posts/users/{Id here}'
* e.g ``` curl.exe http://localhost:10000/posts/users/<id> ```

# Installation and Setup
- All Basics standard Go Language and MongoDB connectivity libraries are installed
- 2 terminal windows(curl or curl.exe) : 
  - ``` go run main.go ```
  - ``` curl.exe https://localhost:10000 ```
- For operations of POST and GET are currently done using command prompt using curl commands

# Features
- Made the ***server thread safe*** using mutex and synchonized the threads to avoid deadlock situations
- Password of the users are securely stored in databased in ***encrypted*** manner using AES encryption method
- For each endpoint respective ***Pagination*** is provided.
- All the necessary and required ***HTTP erros*** will prompt in any case of system failure.
- ***Test cases*** are added to each HTTP request handler.

# Data Types
```
{
	"id"       : some unquie user Id
	"name"     : "name of the insta user"
	"email"    : "email id of insta user"
	"password" : "encrypted password"
}
```
```
{
	"id"              : some unique post id
	"userid"          : unquie Id of respepective user
	"caption"         : "caption of the post uploaded"
	"imageurl"        : "image url of post"
	"postedtimestamp" : "time stamp of posted post"
}
```
