package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FirstDayAtWork/mustracker/config"
	"github.com/FirstDayAtWork/mustracker/controllers"
	"github.com/FirstDayAtWork/mustracker/models"
)

const port int = 2228

func main() {
	// Move to config
	appConfig, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	db := config.MustConnect(
		appConfig.ConnectToDB(),
	)
	// Move this to a Must method?
	err = models.MigrateAccountData(db)
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
	// Routing
	mux.Handle(controllers.LoginPath, app)
	mux.Handle(controllers.RegisterPath, app)
	mux.Handle(controllers.AccountPath, app)
	mux.Handle(controllers.HomePath, app)
	mux.Handle(controllers.AboutPath, app)
	mux.Handle(controllers.DonatePath, app)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
