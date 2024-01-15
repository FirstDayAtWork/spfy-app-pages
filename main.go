package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/entity"
	"github.com/FirstDayAtWork/mustracker/repository"

	"github.com/FirstDayAtWork/mustracker/webserver"
)

const port int = 2228

func main() {
	db, err := repository.Connect(repository.SQLitePath)
	if err != nil {
		panic("ERROR CONNECTING TO DB CANT PROCEED")
	}
	dh := &webserver.DataHandler{DB: db}

	// Migrate schema
	dh.DB.AutoMigrate(&entity.AccountData{})

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/register", dh.HandleRegister)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
