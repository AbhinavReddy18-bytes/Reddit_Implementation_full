package main

import (
	"fmt"
	"log"
	"project4/apis"
	"project4/client"
	"project4/client_rest"
	"project4/engine"
	"project4/performance"
	"sync"
	"time"
)

func simulateClients(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	privateKey, publicKey := client_rest.GenerateKeys()
	username := fmt.Sprintf("user%d", id)

	client_rest.RegisterUser(username, publicKey)

	postContent := fmt.Sprintf("Hello from %s!", username)
	client_rest.CreatePost(username, postContent, privateKey)

	client_rest.FetchPosts()

	client_rest.FetchUserPublicKey(username)

	fmt.Printf("Client %d simulation complete.\n", id)
}

func main() {
	engineInstance := engine.NewEngine()
	go func() {
		fmt.Println("Starting the Engine and REST API Server...")
		apis.StartAPIServer()
	}()

	time.Sleep(2 * time.Second)

	metrics := performance.StartMetrics()
	engineInstance.SetMetrics(metrics)
	log.Println("Initializing Reddit Clone Simulation...")

	privateKeyAlice, publicKeyAlice := client_rest.GenerateKeys()
	privateKeyBob, publicKeyBob := client_rest.GenerateKeys()

	fmt.Println("==== Using REST Client with RSA Signatures ====")
	client_rest.RegisterUser("Alice", publicKeyAlice)
	client_rest.RegisterUser("Bob", publicKeyBob)

	client_rest.CreatePost("Alice", "Hello World! My first post.", privateKeyAlice)
	client_rest.CreatePost("Bob", "Go is awesome!", privateKeyBob)

	client_rest.FetchPosts()
	client_rest.FetchUserPublicKey("Alice")
	client_rest.FetchUserPublicKey("Bob")

	user1 := client.NewClient("user1", engineInstance)
	user2 := client.NewClient("user2", engineInstance)
	user3 := client.NewClient("user3", engineInstance)

	user1.Register()
	user2.Register()
	user3.Register()

	engineInstance.ConnectUser("user1")
	engineInstance.ConnectUser("user2")
	engineInstance.ConnectUser("user3")

	engineInstance.CreateSubreddit("golang")
	engineInstance.CreateSubreddit("technology")
	engineInstance.CreateSubreddit("movies")

	user1.JoinSubreddit("golang")
	user2.JoinSubreddit("golang")
	user3.JoinSubreddit("movies")
	user1.JoinSubreddit("technology")

	user1.LeaveSubreddit("golang")
	user2.LeaveSubreddit("golang")

	postID1, _ := user1.PostInSubreddit("golang", "Hello, Gophers!")
	postID2, _ := user2.PostInSubreddit("golang", "Go is amazing!")
	postID3, _ := user3.PostInSubreddit("movies", "What's your favorite movie?")

	user1.CommentOnPost("golang", postID2, "Absolutely agree!")
	user3.CommentOnPost("movies", postID3, "I love Inception.")

	user2.UpvotePost("golang", postID1)
	user3.DownvotePost("golang", postID2)

	user1.SendMessage("user2", "Hey, have you tried Go modules?")
	user2.SendMessage("user1", "Yes, they're awesome!")

	log.Println("Starting large-scale simulation...")
	client.SimulateClients(engineInstance, 100, 10, 500, 200)

	var wg sync.WaitGroup
	numClients := 10

	fmt.Println("==== Simulating Multiple Clients ====")
	for i := 1; i <= numClients; i++ {
		wg.Add(1)
		go simulateClients(&wg, i)
	}

	wg.Wait()
	fmt.Println("==== Client Simulation Complete ====")

	log.Println("===== USER KARMA =====")
	for _, username := range []string{"user1", "user2", "user3"} {
		karma := engineInstance.ComputeKarma(username)
		log.Printf("Karma for %s: %d\n", username, karma)
	}
	log.Println("======================")

	metrics.Stop()
	metrics.Report()

	log.Println("Simulation and demonstration complete!")
}
