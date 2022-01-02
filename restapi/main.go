package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

type User struct {
	name     string `json:"name"`
	age      string `json:"age"`
	location string `json:"location"`
}

func createClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "gothic-gradient-336615"

	// [END firestore_setup_client_create]
	// Override with -project flags
	flag.StringVar(&projectID, "project", projectID, "The Google Cloud Platform project ID.")
	flag.Parse()

	// [START firestore_setup_client_create]
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

func createUser(w http.ResponseWriter, r *http.Request) {

	// w.Header().Set("Content-Type", "application/json")
	var user User

	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		fmt.Fprintf(w, "Url Param 'name' is missing")
		return
	}

	user.name = name[0]

	age, ok := r.URL.Query()["age"]

	if !ok || len(age[0]) < 1 {
		fmt.Fprintf(w, "Url Param 'age' is missing")
		return
	}

	user.age = age[0]

	location, ok := r.URL.Query()["location"]

	if !ok || len(location[0]) < 1 {
		fmt.Fprintf(w, "Url Param 'location' is missing")
		return
	}

	user.location = location[0]

	ctx := context.Background()
	client := createClient(ctx)
	defer client.Close()

	_, _, err := client.Collection("userbase").Add(ctx, map[string]interface{}{
		"name":     user.name,
		"age":      user.age,
		"location": user.location,
	})
	if err != nil {
		// json.NewEncoder(w).Encode(err)
		fmt.Fprintf(w, "Error")
		return
	}
	fmt.Fprintf(w, "Successfully added "+user.name)
}

func readUser(w http.ResponseWriter, r *http.Request) {

	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		fmt.Fprintf(w, "Url Param 'name' is missing")
		return
	}

	ctx := context.Background()
	client := createClient(ctx)
	defer client.Close()

	iter := client.Collection("userbase").Where("name", "==", name[0]).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Fprintln(w, doc.Data())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to user management application")
}

func main() {
	router := mux.NewRouter()
	// Mock Data

	router.HandleFunc("/api/createUser", createUser).Methods("GET")
	router.HandleFunc("/api/readUser", readUser).Methods("GET")
	router.HandleFunc("/", handler).Methods("GET")

	http.Handle("/", router)

	// log.Fatal(http.ListenAndServe(":8000", router))

	log.Print("starting server...")

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
