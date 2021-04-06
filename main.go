package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Article = []Article{
	{Id: "1", Title: "Golang Tutorial", Desc: "Tutorial for Golang", Content: "Go is good"},
	{Id: "2", Title: "Rest APIs", Desc: "About Rest APIs", Content: "All about REST"},
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// Write a string to the http ResponseWriter
	fmt.Fprintf(w, "Hello homepage!")
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// Read the contents of the request body into reqBody
	reqBody, err := ioutil.ReadAll(r.Body)

	// Error handling. In case of an error while reading the body,
	// send a status of 400 (Bad Request) to the client
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Variable to store the new article
	var art Article

	// Decode/Unmarshal the body json into an Article object
	json.Unmarshal(reqBody, &art)

	// Append the new article to the slice of articles
	Articles = append(Articles, art)

	// Send the new article object to the client
	json.NewEncoder(w).Encode(art)
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	// Encode all articles as JSON and write them to the HTTP ResponseWriter
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	// Get URL variables
	vars := mux.Vars(r)

	// Get the article ID
	id := vars["id"]

	// Iterate over all the articles and send the one with matching ID
	// to the user
	for _, article := range Articles {
		if article.Id == id {
			json.NewEncoder(w).Encode(article)
			return
		}
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	// Get URL variables
	vars := mux.Vars(r)

	// Get the article ID
	id := vars["id"]

	// Read the body of the request
	reqBody, err := ioutil.ReadAll(r.Body)

	// In case of any error, send a 400 (Bad Request) to the client
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Variable to store the updated article data
	var updatedArt Article

	// Parse the body as JSON and unmarshal it to an article
	jerr := json.Unmarshal(reqBody, &updatedArt)

	// In case of an unmarshal error, send a 400 (Bad request) to the client
	if jerr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Iterate over all the articles and look for an article with
	// the same ID as the given ID and update it with the given
	// content
	for i, article := range Articles {
		if article.Id == id {
			art := &Articles[i]
			art.Title = updatedArt.Title
			art.Desc = updatedArt.Desc
			art.Content = updatedArt.Content

			// Send the updated article to the client
			json.NewEncoder(w).Encode(article)
			return
		}
	}

	// If the article was not found, send a 404 (Not Found) to the client
	w.WriteHeader(http.StatusNotFound)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// Get URL variables
	vars := mux.Vars(r)

	// Get the article ID
	id := vars["id"]

	// Iterate over all the articles. If an article with the specified
	// ID is found, remove it from the slice
	for i, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:i], Articles[i+1:]...)
			json.NewEncoder(w).Encode(article)
			return
		}
	}

	// If no article is found, send a status of 404 (Not Found)
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	// Create a new mux Router
	router := mux.NewRouter().StrictSlash(true)

	// Handle root route
	router.HandleFunc("/", homePage)

	// Handle getting all articles
	router.HandleFunc("/articles", returnAllArticles).Methods("GET")

	// Handle posting a new article
	router.HandleFunc("/articles", createNewArticle).Methods("POST")

	// Handle getting single article
	router.HandleFunc("/articles/{id}", returnSingleArticle).Methods("GET")

	// Handle updating an article
	router.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")

	// Handle deleting an article
	router.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9090", router))
}
