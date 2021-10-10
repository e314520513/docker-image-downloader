package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM docker_images ORDER BY id DESC")

	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
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
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

// func Show(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	nId := r.URL.Query().Get("id")
// 	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	emp := Employee{}
// 	for selDB.Next() {
// 		var id int
// 		var name, city string
// 		err = selDB.Scan(&id, &name, &city)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		emp.Id = id
// 		emp.Name = name
// 		emp.City = city
// 	}
// 	tmpl.ExecuteTemplate(w, "Show", emp)
// 	defer db.Close()
// }

// func New(w http.ResponseWriter, r *http.Request) {
// 	tmpl.ExecuteTemplate(w, "New", nil)
// }

// func Edit(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	nId := r.URL.Query().Get("id")
// 	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	emp := Employee{}
// 	for selDB.Next() {
// 		var id int
// 		var name, city string
// 		err = selDB.Scan(&id, &name, &city)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		emp.Id = id
// 		emp.Name = name
// 		emp.City = city
// 	}
// 	tmpl.ExecuteTemplate(w, "Edit", emp)
// 	defer db.Close()
// }

func Download(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		pullImage(name)
		saveImage(name)
		link := "1"
		insForm, err := db.Prepare("INSERT INTO docker_images(name, link) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, link)
		log.Println("INSERT: Name: " + name + " | Link: " + link)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func pullImage(imageName string) {
	cmd := exec.Command("docker", "pull", imageName)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())

	}
	log.Print(string(stdout))
}

func saveImage(imageName string) {
	cmd := exec.Command("docker", "save", "-o", "dockerImages/"+imageName, imageName)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())

	}
	log.Print(string(stdout))
}

// func Update(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	if r.Method == "POST" {
// 		name := r.FormValue("name")
// 		city := r.FormValue("city")
// 		id := r.FormValue("uid")
// 		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		insForm.Exec(name, city, id)
// 		log.Println("UPDATE: Name: " + name + " | City: " + city)
// 	}
// 	defer db.Close()
// 	http.Redirect(w, r, "/", 301)
// }

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM docker_images WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {

	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	// http.HandleFunc("/show", Show)
	// http.HandleFunc("/new", New)
	// http.HandleFunc("/edit", Edit)
	http.HandleFunc("/download", Download)
	// http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
