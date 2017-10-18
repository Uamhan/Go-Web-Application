package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	Message string        //message sent to html template
	Guess   string        //varible holding users guess to be sent to html template
	Result  template.HTML //result message sent to html template
}

type Page struct {
	Title string //varible holds title of page struct
	Body  []byte //holds body of page struct
}

//loads page "title"
//readFile fils body with contents of "title".txt
//returns in a page struct format
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

	var seed = rand.NewSource(time.Now().UnixNano()) //creats a seed for the rand fucntion based on current time.
	var myRand = rand.New(seed)                      //allows us to creat a number based of the seed.
	var randNum int = myRand.Intn(20)                //sets number to be between 1-20

	target := randNum //sets value to be assiigned to cookiie

	var c, err = r.Cookie("Target") //checks to see if cookie with a name of Target exists and if so sets it to c
	//if cookiie is set
	if err == nil {
		//we set target to the value of target in the cookie
		target, _ = strconv.Atoi(c.Value)
	}
	//we create the cookie
	cookie := &http.Cookie{Name: "Target", Value: strconv.Itoa(target)} //creates a new http.cookie and assigns values to it
	//set the cookie
	http.SetCookie(w, cookie)

	r.ParseForm()                                  //parses form so we can use values from it
	m := Message{}                                 //creats blank Message struct
	guess, _ := strconv.Atoi(r.FormValue("Guess")) //converts the Guess input from the from to an int
	//target, _ := strconv.ParseInt(c.Value, 10, 0)             // converts the target value from cookie to an int

	//if no guesses have been made yet
	if r.FormValue("Guess") == "" {
		m = Message{Message: "Guess a number between 1 and 20 "}
		//guess has been made and the guess is an integer
	} else if _, err := strconv.Atoi(r.FormValue("Guess")); err == nil {

		//guess and targert are equal
		if guess == target {
			m = Message{Message: "Guess a number between 1 and 20 ", Guess: "Your Guess is : " + r.FormValue("Guess"), Result: "Congratulations you Guessed Correctly </br> <a href=\"/Guess\">New Game</a> "}
			//resets cookie after correct guess
			target = randNum
			cookie = &http.Cookie{Name: "Target", Value: strconv.Itoa(target)}
			http.SetCookie(w, cookie)
			//c, _ = r.Cookie("Target")
			//guess is greater than target
		} else if guess > target {
			m = Message{Message: "Guess a number between 1 and 20 ", Guess: "Your Guess is : " + r.FormValue("Guess"), Result: "You Guessed to High"}
			//guess is less than target
		} else {
			m = Message{Message: "Guess a number between 1 and 20 ", Guess: "Your Guess is " + r.FormValue("Guess"), Result: "You Guessed to Low"}
		}

		//invalid guess has been made
	} else {
		m = Message{Message: "Guess a number between 1 and 20 ", Guess: "Invalid Guess"}
	}
	t, _ := template.ParseFiles("Guess.tmpl") //parses guess.tmpl into a template
	t.Execute(w, m)                           // executes the template we just created and fills it with Message struct created above

}

//handler for home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("index")
	fmt.Fprintf(w, "<div class = \"container text-center\"> </br></br><h1>%s</h1><p>%s</p></div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/", homeHandler)       //tells the http package to handle all requests to the web root
	http.HandleFunc("/Guess", GuessHandler) //tells the http package to handle all requests to /Guess
	http.ListenAndServe(":8080", nil)       //listens and servers requests on port :8080
}
