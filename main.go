package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Mess struct {
	Id    string `json:"Id"`
	Title string `json:"Title"`
	Desc  string `json:"desc"`
}

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}
type Bank struct {
	Id                  uint16
	Name                string
	Loan, Percent, Term float64
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
var postsBank = []Bank{}
var showPostBanks = Bank{}
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
func banks(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/bank.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db := dbConn()

	//виборка даних
	res, err := db.Query("SELECT  *   FROM  `banks` ")
	if err != nil {
		panic(err)
	}

	postsBank = []Bank{}

	for res.Next() {
		var post Bank
		err = res.Scan(&post.Id, &post.Name, &post.Loan, &post.Percent, &post.Term)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Post: %s with id %d", post.Name, post.Id))

		postsBank = append(postsBank, post)
	}

	t.ExecuteTemplate(w, "bank", postsBank)
	defer db.Close()
}
func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}
func createBank(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/createBank.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "createBank", nil)
}
func elif(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/elif.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "elif", nil)
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
func saveBank(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	loan := r.FormValue("loan")
	percent := r.FormValue("percent")
	term := r.FormValue("term")

	if name == "" || loan == "" || percent == "" || term == "" {
		fmt.Fprintf(w, "Не всі поля заповнені")
	} else {

		db := dbConn()

		//встановлення даних
		insert, err := db.Query(fmt.Sprintf("INSERT  INTO `banks`(`name`, `loan`,`percent`,`term`) VALUES ('%s', '%s','%s','%s')", name, loan, percent, term))

		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/banks", http.StatusSeeOther)

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
func showBank(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/showBank.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db := dbConn()

	nId := r.URL.Query().Get("id")

	//виборка даних
	res, err := db.Query(fmt.Sprintf("SELECT  *   FROM  `banks` WHERE `id` = '%s' ", nId))
	if err != nil {
		panic(err)
	}

	showPostBanks = Bank{}

	for res.Next() {
		var postsBank Bank
		err = res.Scan(&postsBank.Id, &postsBank.Name, &postsBank.Loan, &postsBank.Percent, &postsBank.Term)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Post: %s with id %d", postsBank.Name, postsBank.Id))

		showPostBanks = postsBank
	}

	t.ExecuteTemplate(w, "showBank", showPostBanks)
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
func editBank(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/editBank.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM banks WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	art := Bank{}
	for selDB.Next() {
		var id uint16
		var name string
		var loan, percent, term float64
		err = selDB.Scan(&id, &name, &loan, &percent, &term)
		if err != nil {
			panic(err.Error())
		}
		art.Id = id
		art.Name = name
		art.Loan = loan
		art.Percent = percent
		art.Term = term
	}
	t.ExecuteTemplate(w, "editBank", art)
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
func updateBank(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		loan := r.FormValue("loan")
		percent := r.FormValue("percent")
		term := r.FormValue("term")

		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE banks SET name=?, loan=?, percent=?, term=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, loan, percent, term, id)
		log.Println("UPDATE: name: " + name + " | loan: " + loan + " | percent: " + percent + " | terms" + term)
	}
	defer db.Close()
	http.Redirect(w, r, "/banks", 301)
}
func calculate(w http.ResponseWriter, r *http.Request) {
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

/* Ajax */
func calFormula(p float64, r float64, n float64) (m float64) {
	var res0, res1, res2 float64

	res0 = r / 12
	res1 = math.Pow(1+res0, n)
	res2 = p * res0
	m = res2 * res1 / (res1 - 1)

	return m
}
func receiveAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		val1 := r.FormValue("val1")
		val2 := r.FormValue("val2")
		val3 := r.FormValue("val3")

		strconv.ParseFloat(val1, 64)
		strconv.ParseFloat(val2, 64)
		strconv.ParseFloat(val3, 64)

		result1, _ := strconv.ParseFloat(val1, 64)
		result2, _ := strconv.ParseFloat(val2, 64)
		result3, _ := strconv.ParseFloat(val3, 64)

		result := calFormula(result1, result2, result3)

		fmt.Println("P=", result1)
		fmt.Println("r=", result2)
		fmt.Println("n=", result3)

		resultString := strconv.FormatFloat(result, 'E', -1, 64)

		fmt.Println("Sum: ", result)
		w.Write([]byte(resultString))
	}
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
func deleteBank(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	art := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM banks WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(art)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/banks", 301)
}
func confirmation(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/confirmation.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "confirmation", nil)
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

/* API */
func getArticles(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/art.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	resp, err := http.Get("http://192.168.1.18/api/articles")
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
func handleFunc() {
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/banks", banks).Methods("GET")
	rtr.HandleFunc("/delete", DeleteA)
	rtr.HandleFunc("/deleteBank", deleteBank)
	rtr.HandleFunc("/update", UpdateA)
	rtr.HandleFunc("/updateBank", updateBank)
	rtr.HandleFunc("/calculate", calculate)
	rtr.HandleFunc("/edit", EditA)
	rtr.HandleFunc("/editBank", editBank)
	rtr.HandleFunc("/receive", receiveAjax)
	rtr.HandleFunc("/count", receiveAjax)
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/createBank", createBank).Methods("GET")
	rtr.HandleFunc("/elif", elif).Methods("GET")
	rtr.HandleFunc("/art", getArticles).Methods("GET")
	rtr.HandleFunc("/contacts", contacts).Methods("GET")
	rtr.HandleFunc("/confirmation", confirmation).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/saveBank", saveBank).Methods("POST")
	rtr.HandleFunc("/post", show_post).Methods("GET")
	rtr.HandleFunc("/bank", showBank).Methods("GET")
	http.Handle("/", rtr)
	log.Println("Server started on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	//http.ListenAndServe("192.168.1.9:1111", nil)
}

func main() {

	handleFunc()

}

//qq
