# Instagram-Backend-API(NO EXTRA DEPENDENCIES)

Designed and Developed an HTTP JSON REST API mocking of instagram posts and users schema which capable of the following structure operations given below. 
Golang is used for the API and MongoDB is used as a storage.

# Structure
- [x] Create an User
* Should be a POST request
* Use JSON request body
* URL should be ‘/users'

- [x] Get a user using id
* Should be a GET request
* Id should be in the url parameter
* URL should be ‘/users/{id here}’
  
- [x] Create a Post
* Should be a POST request
* Use JSON request body
* URL should be ‘/posts'
  
- [x] Get a post using id
* Should be a GET request
* Id should be in the url parameter
* URL should be ‘/posts/{id here}
  
- [x] List all posts of a perticular user
* Should be a GET request
* URL should be ‘/posts/users/{Id here}'

# Installation and Setup
- All Basics standard Go Language and MongoDB connectivity libraries are installed

# Features
- Made the server thread safe using mutex and synchonized the threads to avoid deadlock situations
- Password of the users are securely stored in databased in encrypted manner usinf AES encryption method
- For each endpoint respective Pagination is provided.
- 


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
