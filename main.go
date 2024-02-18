package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/controllers"
	"github.com/FirstDayAtWork/mustracker/models"
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
	err = models.MigrateAccessTokenData(db)
	if err != nil {
		panic(err)
	}

	r := &models.Repository{
		DB: db,
	}
	auth := &controllers.Authorizer{
		Secret:     "some-secret-wooow",
		Issuer:     "application",
		Repository: r,
	}
	th := controllers.NewTemplateHandler()
	app := &controllers.App{
		Th:         th,
		Repository: r,
		Auth:       auth,
	}

	// Server boilerplate
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// Parse template - configure handler - add to router
	mux.Handle(controllers.LoginPath, app)
	mux.Handle(controllers.RegisterPath, app)
	mux.Handle(controllers.AccountPath, app)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
