package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}
var db *sql.DB
var err error

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//виборка даних
	res, err := db.Query("SELECT  *   FROM  `articles` ")
	if err != nil {
		panic(err)
	}

	posts = []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Post: %s with id %d", post.Title, post.Id))

		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}
func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}
func contacts(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/contacts.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "contacts", nil)
}
func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не всі поля заповнені")
	} else {

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		//встановлення даних
		insert, err := db.Query(fmt.Sprintf("INSERT  INTO `articles`(`title`, `anons`,`full_text`) VALUES ('%s', '%s','%s')", title, anons, full_text))

		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}
func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//виборка даних
	res, err := db.Query(fmt.Sprintf("SELECT  *   FROM  `articles` WHERE `id` = '%s' ", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Post: %s with id %d", post.Title, post.Id))

		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)
}
func send(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate form
	msg := &Message{
		Email:   r.PostFormValue("email"),
		Content: r.PostFormValue("content"),
	}

	if msg.Validate() == false {
		render(w, "templates/contacts.html", msg)
		return
	}

	// Step 2: Send contact form message in an email
	if err := msg.Deliver(); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	// Step 3: Redirect to confirmation page
	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}
func confirmation(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/confirmation.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "confirmation", nil)
}
func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
func loginG(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "login", nil)
}
func signupG(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/signup.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "signup", nil)
}
func loginP(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(w, r, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login", 301)
		return
	}

	w.Write([]byte("Hello " + databaseUsername))
}
func signupP(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}

		w.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(w, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(w, r, "/", 301)
	}
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/login", loginG).Methods("GET")
	rtr.HandleFunc("/signup", signupG).Methods("GET")
	rtr.HandleFunc("/contacts", contacts).Methods("GET")
	rtr.HandleFunc("/contacts", send).Methods("POST")
	rtr.HandleFunc("/confirmation", confirmation).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)

	//http.HandleFunc("/signup", signupPage)
	//http.HandleFunc("/login", loginPage)
	//rtr.HandleFunc("/login", loginG).Methods("GET")
	//rtr.HandleFunc("/signup", signupG).Methods("GET")

	rtr.HandleFunc("/login", loginG).Methods("GET")
	rtr.HandleFunc("/signup", signupG).Methods("GET")
	//rtr.HandleFunc("/login", loginP).Methods("POST")
	//rtr.HandleFunc("/signup", signupP).Methods("POST")

}

func main() {
	handleFunc()
}
