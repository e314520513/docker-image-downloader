package main

import (
	"bytes"
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const imageDir = "dockerImages/"
const imageExt = ".tar.gz"

type Images struct {
	Id   int
	Name string
	Link string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbName := "golang"
        dbPassword :="golang"
	protocal := "tcp(localhost:3306)"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@"+protocal+"/"+dbName)
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
	emp := Images{}
	res := []Images{}
	for selDB.Next() {
		var id int
		var name, link string
		err = selDB.Scan(&id, &name, &link)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Link = link
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
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM docker_images where id=?", nId)

	if err != nil {
		panic(err.Error())
	}
	var id int
	var name, link string
	for selDB.Next() {

		err = selDB.Scan(&id, &name, &link)
		if err != nil {
			panic(err.Error())
		}
		log.Println(link)
	}

	// 讀取檔案
	downloadBytes, err := ioutil.ReadFile(link)

	if err != nil {
		log.Println(err)
	}

	// 取得檔案的 MIME type
	mime := http.DetectContentType(downloadBytes)

	fileSize := len(string(downloadBytes))

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))

	http.ServeContent(w, r, link, time.Now(), bytes.NewReader(downloadBytes))
	log.Println("downloaded")
	defer db.Close()
	http.Redirect(w, r, "/", 301)

}
func Search(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		pullImage(name)
		imagePath, exported := saveImage(name)

		_, err := os.Stat(imagePath)

		if err == nil && exported {
			insForm, err := db.Prepare("INSERT INTO docker_images(name, link) VALUES(?,?)")
			if err != nil {
				panic(err.Error())
			}

			if _, execErr := insForm.Exec(name, imagePath);execErr != nil{
				panic(execErr.Error())
			}
			log.Println("pulled it! Name: " + name + " | Link: " + imagePath)
		}

	}
	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func pullImage(name string) {
	log.Println("pulling image")
	cmd := exec.Command("/bin/bash", "-c", "docker pull "+name)
	execCMD(cmd)

}

func saveImage(name string) (string, bool) {
	var exported = false
	log.Println("saving image")

	imagePath := imageDir + strings.Replace(name, "/", "_", -1) + imageExt

	fileinfo, _ := os.Stat(imagePath)

	if fileinfo == nil {
		if err := os.Mkdir(imageDir, 0755); err != nil {
			log.Println(err)
		}

		cmd := exec.Command("/bin/bash", "-c", "docker save -o "+imagePath+" "+name)
		execCMD(cmd)
		cmd = exec.Command("/bin/bash", "-c", "tar -zcvf "+imagePath+" "+imagePath)
		execCMD(cmd)
		exported = true
	}

	return imagePath, exported
}

func execCMD(cmd *exec.Cmd) {
	//建立獲取命令輸出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	//執行命令
	if err := cmd.Start(); err != nil {
		log.Println("Error:The command is err,", err)
		return
	}
	//讀取所有輸出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ReadAll Stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err.Error())
		return
	}
	log.Printf("stdout:\n\n %s", bytes)
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
	imagePath := r.URL.Query().Get("link")

	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM docker_images WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	os.Remove(imagePath)
	log.Println("Deleted")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
