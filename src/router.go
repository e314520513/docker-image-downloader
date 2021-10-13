package main

import (
	"log"
	"net/http"
)

func router() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/download", Download)
	// http.HandleFunc("/show", Show)
	// http.HandleFunc("/new", New)
	// http.HandleFunc("/edit", Edit)
	http.HandleFunc("/search", Search)
	// http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
