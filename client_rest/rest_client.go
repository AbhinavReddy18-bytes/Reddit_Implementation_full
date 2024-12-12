package client_rest

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type Post struct {
	User      string `json:"user"`
	Content   string `json:"content"`
	Signature string `json:"signature"`
}

func GenerateKeys() (*rsa.PrivateKey, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return nil, ""
	}

	publicKeyBytes := base64.StdEncoding.EncodeToString(privateKey.PublicKey.N.Bytes())
	return privateKey, publicKeyBytes
}

func SignMessage(privateKey *rsa.PrivateKey, message string) string {
	hash := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hash[:])
	if err != nil {
		fmt.Println("Error signing message:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

func RegisterUser(username, publicKey string) {
	user := User{Username: username, PublicKey: publicKey}
	payload, _ := json.Marshal(user)

	resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error registering user:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Register Response:", string(body))
}

func CreatePost(username, content string, privateKey *rsa.PrivateKey) {
	signature := SignMessage(privateKey, content)
	post := Post{User: username, Content: content, Signature: signature}

	payload, _ := json.Marshal(post)
	resp, err := http.Post("http://localhost:8080/post", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating post:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Create Post Response:", string(body))
}

func FetchPosts() {
	resp, err := http.Get("http://localhost:8080/posts")
	if err != nil {
		fmt.Println("Error fetching posts:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("All Posts:", string(body))
}

func FetchUserPublicKey(username string) {
	url := fmt.Sprintf("http://localhost:8080/user/%s/publickey", username)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching public key:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("User Public Key:", string(body))
}
