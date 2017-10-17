package main

import (
    "fmt"	
	"io/ioutil"	
    "net/http"		
	"html/template"	
)

type Message struct {
	Message string
}

type Page struct {
    Title string
    Body  []byte
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: "Guessing Game", Body: body}, nil
}

//function handler takes a ResponseWriter and a Request asargument
//http.ResponseWriter assembles Http servers response when written to
//http.Request represents the clients http request
func Handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
    p, _ := loadPage(title)
	//if no page is loaded loads home page
	if p == nil{
		p, _ = loadPage("index")
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	}else{
	m := Message{Message: "Guess a number between 1 and 20"}
	t, _ := template.ParseFiles("Guess.tmpl")
	t.Execute(w,m)
	}
}

//func homeHandler(w http.ResponseWriter, r *http.Request) {
 //   p, _ := loadPage("index")
 //   fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)	
//	}


func main() {
    http.HandleFunc("/", Handler)  //tells the http package to handle all requests to the web root 
    http.ListenAndServe(":8080", nil) //listens and servers requests on port :8080
}