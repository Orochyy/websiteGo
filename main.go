package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

//type Mess struct {
//	name string `json:"name"`
//	desc string `json:"desc"`
//}

type Mess struct {
	Id    string `json:"Id"`
	Title string `json:"Title"`
	Desc  string `json:"desc"`
}

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}
type Employee struct {
	Id   int
	Name string
	City string
}

type Articleq struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Articleq

var mess []Mess
var tmpl = template.Must(template.ParseGlob("templates/*"))
var posts = []Article{}
var showPost = Article{}
var db *sql.DB
var err error

func dbConn() (db *sql.DB) {

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db

}
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db := dbConn()

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
	defer db.Close()
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

		db := dbConn()

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

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db := dbConn()

	nId := r.URL.Query().Get("id")

	//виборка даних
	res, err := db.Query(fmt.Sprintf("SELECT  *   FROM  `articles` WHERE `id` = '%s' ", nId))
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
	defer db.Close()
}
func ShowE(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "show2", emp)
	defer db.Close()
}
func EditA(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/edit.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM articles WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	art := Article{}
	for selDB.Next() {
		var id uint16
		var title, anons, full_text string
		err = selDB.Scan(&id, &title, &anons, &full_text)
		if err != nil {
			panic(err.Error())
		}
		art.Id = id
		art.Title = title
		art.Anons = anons
		art.FullText = full_text
	}
	t.ExecuteTemplate(w, "edit", art)
	defer db.Close()
}
func UpdateA(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		title := r.FormValue("title")
		anons := r.FormValue("anons")

		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE articles SET title=?, anons=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(title, anons, id)
		log.Println("UPDATE: title: " + title + " | anons: " + anons)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
func DeleteA(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	art := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM articles WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(art)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
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
func signupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/signup.html")
		return
	}

	db := dbConn()

	defer db.Close()

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user string

	err = db.QueryRow("SELECT username FROM usersl WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO usersl(username, password) VALUES(?, ?)", username, hashedPassword)
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
func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	db := dbConn()

	defer db.Close()

	username := r.FormValue("username")
	password := r.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err = db.QueryRow("SELECT username, password FROM usersl WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}

	t, err := template.ParseFiles("templates/indexU.html", "templates/headerU.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "indexU", posts)

	//w.Write([]byte("Hello " + databaseUsername))

}
func getArticles(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/art.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	resp, err := http.Get("http://192.168.1.6/api/articles")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
	t.ExecuteTemplate(w, "art", posts)
}
func vlaDick(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(mess)

}
func handleFunc() {
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/delete", DeleteA)
	rtr.HandleFunc("/update", UpdateA)
	rtr.HandleFunc("/edit", EditA)
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/art", getArticles).Methods("GET")
	rtr.HandleFunc("/contacts", contacts).Methods("GET")
	rtr.HandleFunc("/contacts", send).Methods("POST")
	rtr.HandleFunc("/confirmation", confirmation).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post", show_post).Methods("GET")
	rtr.HandleFunc("/dick", vlaDick).Methods("GET")
	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Println("Server started on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	//http.ListenAndServe("192.168.1.9:80", nil)

}

func main() {
	mess = []Mess{
		Mess{Id: "VlaDick", Title: "loh", Desc: "double loh^^ kek "},
	}

	handleFunc()

}

//qq
