package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"

	"github.com/FirstDayAtWork/mustracker/controllers"
)

const port int = 2228

func main() {
	db, err := models.Connect(models.SQLitePath)
	if err != nil {
		panic("ERROR CONNECTING TO DB CANT PROCEED")
	}
	dh := &controllers.DataHandler{DB: db}

	// Migrate schema
	dh.DB.AutoMigrate(&models.AccountData{})

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/register", dh.HandleRegister)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
