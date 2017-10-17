package main

import (
    "fmt"	
	"io/ioutil"	
    "net/http"		
	"html/template"	
	"math/rand"
	"time"
	"strconv"

)

type Message struct {
	Message string
	Guess string
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
func GuessHandler(w http.ResponseWriter, r *http.Request) {
	
	var seed = rand.NewSource(time.Now().UnixNano()) 	//creats a seed for the rand fucntion based on current time.
	var myRand = rand.New(seed)							//allows us to creat a number based of the seed.
	var randNum int = myRand.Intn(20)
	
	cookie := http.Cookie{ Name:  "Target", Value: fmt.Sprint(randNum)}
	
	c, err := r.Cookie("Target")
	
	if err != nil {
		http.SetCookie(w, &cookie) //sets cookie if there isint one already
		c, _ = r.Cookie("Target")
	}
	fmt.Println(c)
	if r.Method == "GET" {
		
        r.ParseForm()
        fmt.Println("Guess:", r.Form["Guess"])
		m := Message{}
		if r.FormValue("Guess") == ""{
		m = Message{Message: "Guess a number between 1 and 20 " }	
		}else if _, err := strconv.Atoi(r.FormValue("Guess")); err == nil {
			m = Message{Message: "Guess a number between 1 and 20 ", Guess: "Your Guess is : " + r.FormValue("Guess") }
		}else{
			m = Message{Message: "Guess a number between 1 and 20 ", Guess: "Invalid Guess" }	
		}
		t, _ := template.ParseFiles("Guess.tmpl")
		t.Execute(w,m)
    }
	
	
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
   p, _ := loadPage("index")
   fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)	
}


func main() {
	http.HandleFunc("/",homeHandler) //tells the http package to handle all requests to the web root
    http.HandleFunc("/Guess", GuessHandler)  //tells the http package to handle all requests to /Guess 
    http.ListenAndServe(":8080", nil) //listens and servers requests on port :8080
}