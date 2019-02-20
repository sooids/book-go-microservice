package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldResponse struct {
	Message string `json:"message"`
	Author  string `json:"-"`
	Date    string `json:",omitempty"`
	ID      int    `json:"id,string"`
}

type helloWorldRequest struct {
	Name string `json:"name"`
}

type validationHandler struct {
	next http.Handler
}

type helloWorldHandler struct{}

func main() {
	port := 8080
	http.HandleFunc("/helloworld", helloWorldHandleFunc)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	log.Printf("Server starting on port: %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandleFunc(w http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	response := helloWorldResponse{Message: "Hello " + request.Name}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(&response); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}
}

func newValidationHandler(next http.Handler) http.Handler {
	return validationHandler{next: next}
}

func newHelloWorldHandler() http.Handler {
	return helloWorldHandler{}
}

func (h validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&request); err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}

	h.next.ServeHTTP(rw, r)
}

func (h helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	response := helloWorldResponse{Message: "Hello"}
	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(response); err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		log.Printf(err.Error())
		return
	}
}
