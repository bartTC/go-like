package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	counter       = make(map[string]int)
	statsFilePath = "/tmp/like.json"
)

func logIf(err error) {
	if err != nil {
		log.Println(err)
	}
}

func writeCounterData() {
	f, err := os.Create(statsFilePath)
	logIf(err)

	defer f.Close()

	data, err := json.Marshal(counter)
	logIf(err)

	_, err = f.Write(data)
	logIf(err)
}

func buttonHandler(w http.ResponseWriter, r *http.Request) {
	object := r.URL.Path[1:]
	method := r.Method

	type Context struct {
		Object  string `json:"object"`
		Counter int    `json:"counter"`
	}

	// The first and only URL path defines the object name we use as an
	// identifier to retrieve the counter.
	if len(object) == 0 {
		http.Error(w, "No object given", http.StatusBadRequest)
		return
	}

	// In case this object is not yet registered, just add it to the list
	if _, ok := counter[object]; !ok {
		counter[object] = 0
	}

	// POST requests increae the counter by one
	if method == "POST" {
		counter[object] += 1
	}

	// Return the current counter for the given object name
	context := &Context{Object: object, Counter: counter[object]}
	response, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(response))
	writeCounterData()
}

func main() {
	log.Println("Starting server on port 8080...")

	content, err := ioutil.ReadFile(statsFilePath)
	logIf(err)

	err = json.Unmarshal(content, &counter)
	logIf(err)

	http.HandleFunc("/", buttonHandler)
	http.ListenAndServe(":8081", nil)
}

