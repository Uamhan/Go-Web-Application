package main

import (
    "fmt"				
    "net/http"			
)
//function handler takes a ResponseWriter and a Request asargument
//http.ResponseWriter assembles Http servers response when written to
//http.Request represents the clients http request
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Guessing Game</h1>")
}

func main() {
    http.HandleFunc("/", handler)  //tells the http package to handle all requests to the web root 
    http.ListenAndServe(":8080", nil) //listens and servers requests on port :8080
}