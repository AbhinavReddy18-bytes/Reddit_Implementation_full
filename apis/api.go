package apis

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Post struct {
	User      string `json:"user"`
	Content   string `json:"content"`
	Signature string `json:"signature"`
}

type User struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

var (
	posts []Post
	users = make(map[string]string)
	mu    sync.Mutex
)

func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	users[user.Username] = user.PublicKey
	mu.Unlock()

	response := map[string]string{"status": "registered", "user": user.Username}
	json.NewEncoder(w).Encode(response)
}

func verifySignature(content, signature, publicKey string) bool {
	pubKeyBytes, _ := base64.StdEncoding.DecodeString(publicKey)
	hash := sha256.Sum256([]byte(content))
	sigBytes, _ := base64.StdEncoding.DecodeString(signature)

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(pubKeyBytes),
		E: 65537,
	}

	err := rsa.VerifyPKCS1v15(pubKey, 0, hash[:], sigBytes)
	return err == nil
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	publicKey, exists := users[post.User]
	mu.Unlock()

	if !exists || !verifySignature(post.Content, post.Signature, publicKey) {
		http.Error(w, "Signature verification failed", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	posts = append(posts, post)
	mu.Unlock()

	response := map[string]string{"status": "post created", "user": post.User}
	json.NewEncoder(w).Encode(response)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(posts)
}

func getUserPublicKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	mu.Lock()
	publicKey, exists := users[username]
	mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]string{"username": username, "public_key": publicKey}
	json.NewEncoder(w).Encode(response)
}

func StartAPIServer() {
	router := mux.NewRouter()
	router.HandleFunc("/register", registerUser).Methods("POST")
	router.HandleFunc("/post", createPost).Methods("POST")
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/user/{username}/publickey", getUserPublicKey).Methods("GET")

	fmt.Println("API Server is running on port :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
