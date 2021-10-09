package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id       int    `json:"id" bson:"id,omitempty"`
	Name     string `json:"name" bson:"name,omitempty"`
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}

type Post struct {
	Id              int       `json:"id" bson:"id,omitempty"`
	Userid          int       `json:"userid" bson:"userid,omitempty"`
	Caption         string    `json:"caption" bson:"caption,omitempty"`
	ImageURL        string    `json:"imageurl" bson:"imageurl,omitempty"`
	PostedTimestamp time.Time `json:"postedtimestamp" bson:"postedtimestamp,omitempty"`
}

type userHandlers struct {
	sync.Mutex
	store map[int]User
}

type postHandlers struct {
	sync.Mutex
	store map[int]Post
}

func (u *userHandlers) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u.getusers(w, r)
		return
	case "POST":
		u.postusers(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (p *postHandlers) posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p.getposts(w, r)
		return
	case "POST":
		p.postposts(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

// ConnectDB function
func ConnectDB(db string, c string) *mongo.Collection {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection := client.Database(db).Collection(c)
	return collection
}

// ErrorResponse
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError function
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

func (u *userHandlers) getusers(w http.ResponseWriter, r *http.Request) {
	var usercollection *mongo.Collection = ConnectDB("instadb", "user")
	var users []User
	cur, err := usercollection.Find(context.TODO(), bson.M{})
	if err != nil {
		GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())
	u.Lock()
	for cur.Next(context.TODO()) {
		var user User
		// & character returns the memory address of the following variable.
		err := cur.Decode(&user) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	u.Unlock()
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(users) // encode similar to serialize process.

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint Hit: returnAllUsers")
}

func (p *postHandlers) getposts(w http.ResponseWriter, r *http.Request) {
	var postcollection *mongo.Collection = ConnectDB("instadb", "post")
	var posts []Post
	cur, err := postcollection.Find(context.TODO(), bson.M{})
	if err != nil {
		GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())
	p.Lock()
	for cur.Next(context.TODO()) {
		var post Post
		err := cur.Decode(&post)
		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, post)
	}
	p.Unlock()

	json.NewEncoder(w).Encode(posts)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint Hit: returnAllPosts")
}

// Hash key is sent specifically for encoding key
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// AES encrption the password and stored in database
func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

//not yet used since user not logging or any thing else so to show the password now
func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

// Create an User
// Should be a POST request
// Use JSON request body
// URL should be ‘/users'

func (u *userHandlers) postusers(w http.ResponseWriter, r *http.Request) {
	var usercollection *mongo.Collection = ConnectDB("instadb", "user")

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'\n", ct)))
		return
	}

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Password = string(encrypt([]byte(user.Password), "locking"))
	result, err := usercollection.InsertOne(context.TODO(), user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("Endpoint Hit: GettinguserJSON")

	u.Lock()
	json.NewEncoder(w).Encode(result)
	defer u.Unlock()
}

// Create a Post
// Should be a POST request
// Use JSON request body
// URL should be ‘/posts'

func (p *postHandlers) postposts(w http.ResponseWriter, r *http.Request) {
	var postcollection *mongo.Collection = ConnectDB("instadb", "post")

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	post.PostedTimestamp = time.Now()
	result, err := postcollection.InsertOne(context.TODO(), post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("Endpoint Hit: GettingpostJSON")

	p.Lock()
	json.NewEncoder(w).Encode(result)
	defer p.Unlock()
}

// Get a user using id
// Should be a GET request
// Id should be in the url parameter
// URL should be ‘/users/<id here>’

func (u *userHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var usercollection *mongo.Collection = ConnectDB("instadb", "user")

	u.Lock()
	val, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var user User
	filter := bson.M{"id": val}
	error := usercollection.FindOne(context.TODO(), filter).Decode(&user)
	u.Unlock()
	if error != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(user)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint Hit: Gettingperticularuser")
}

// Get a post using id
// Should be a GET request
// Id should be in the url parameter
// URL should be ‘/posts/<id here>’

func (p *postHandlers) getPost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var postcollection *mongo.Collection = ConnectDB("instadb", "post")
	p.Lock()
	val, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var post Post
	filter := bson.M{"id": val}
	error := postcollection.FindOne(context.TODO(), filter).Decode(&post)
	p.Unlock()
	if error != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(post)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint Hit: Gettingperticularpost")
}

// List all posts of a user
// Should be a GET request
// URL should be ‘/posts/users/<Id here>'

func (p *postHandlers) allPosts(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	val, error := strconv.Atoi(parts[3])
	if error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var postcollection *mongo.Collection = ConnectDB("instadb", "post")
	var posts []Post
	cur, err := postcollection.Find(context.TODO(), bson.M{})
	if err != nil {
		GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())
	p.Lock()
	for cur.Next(context.TODO()) {

		var post Post
		err := cur.Decode(&post) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		if post.Userid == val {
			posts = append(posts, post)
		}
	}
	p.Unlock()

	json.NewEncoder(w).Encode(posts)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint Hit: GettingallpostofPerticularuser")
}

func newPostHandlers() *postHandlers {
	return &postHandlers{
		store: map[int]Post{},
	}
}

func newUserHandlers() *userHandlers {
	return &userHandlers{
		store: map[int]User{},
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: homePage")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("<html><h1><b>Welcome to the Instagram Backend API</b></h1></html>"))
}

func main() {
	http.HandleFunc("/", homePage)
	userHandlers := newUserHandlers()
	postHandlers := newPostHandlers()
	http.HandleFunc("/users", userHandlers.users)
	http.HandleFunc("/posts", postHandlers.posts)
	http.HandleFunc("/users/", userHandlers.getUser)
	http.HandleFunc("/posts/", postHandlers.getPost)
	http.HandleFunc("/posts/users/", postHandlers.allPosts)
	log.Println("Listening on localhost:10000")
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		panic(err)
	}
}
