package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/FirstDayAtWork/mustracker/controllers"
	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/FirstDayAtWork/mustracker/templates"
	"github.com/FirstDayAtWork/mustracker/views"
)

const port int = 2228

func main() {
	// Move to config
	DBconfig := &models.SQLiteConfig{
		StorageDir:  "local_storage",
		Environment: "local_dev",
		DBName:      "db",
	}
	db := models.MustConnect(
		DBconfig.ConnectToDB(),
	)
	// Move this to a Must method?
	err := models.MigrateAccountData(db)
	if err != nil {
		panic(err)
	}

	r := &models.Repository{
		DB: db,
	}

	// Server boilerplate
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// Parse template - configure handler - add to router
	regTpl := views.Must(
		views.ParseFS(
			templates.FS,
			filepath.Join("base.html"),
			filepath.Join("register.html"),
		),
	)
	mux.Handle(
		"/register",
		&controllers.RegisterHandler{
			Tpl:        regTpl,
			Repository: r,
		},
	)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
