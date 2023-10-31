package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)

type Kitty struct {
	Id int
	Name string
	Age int
	IsSecretlyEvil bool
}

const PORT = 8888

func getDbConnection() *sql.DB{
	db, err := sql.Open("sqlite3", "./kittens.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Kitty")
	fmt.Fprintf(w, "Hello Kitty")
}

func kittenHandler(w http.ResponseWriter, r *http.Request) {
	switch(r.Method) {
		case "GET":
			getKittens(w, r)
		case "POST":
			createKitten(w,r)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func getKittens(w http.ResponseWriter, r *http.Request) {
	db := getDbConnection()
	stmt := "SELECT * FROM kittens"
	rows,err := db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var scanner Kitty
	var kittens []Kitty
	for rows.Next() {
		_ = rows.Scan(&scanner.Id, &scanner.Name, &scanner.Age, &scanner.IsSecretlyEvil)
		kittens = append(kittens, scanner)
	}
	v, err := json.Marshal(kittens)
	if err != nil {
		fmt.Fprintf(w, "Rip")
	}
	fmt.Fprintf(w, "%s", string(v))
}

func createKitten(w http.ResponseWriter, r *http.Request){
	var kitty Kitty
	err := json.NewDecoder(r.Body).Decode(&kitty)
	if err != nil {
		w.WriteHeader(400)
		res, err := json.Marshal("Bad request bitch.") 
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(res)
		w.Write(res)
	}

	db := getDbConnection()
	stmt, err := db.Prepare("INSERT INTO kittens(name, age, isSecretlyEvil) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(kitty.Name, kitty.Age, kitty.IsSecretlyEvil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%+v", kitty)
}

func main() {
	db:= getDbConnection()
	stmt := `
		CREATE TABLE IF NOT EXISTS kittens (id integer not null primary key autoincrement, name text, age int, isSecretlyEvil boolean)
	`
	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/",index)
	http.HandleFunc("/kittens",kittenHandler)
	http.ListenAndServe(":8080", nil)
}
