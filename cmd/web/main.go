package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RistoFlink/bookings/internal/config"
	"github.com/RistoFlink/bookings/internal/handlers"
	"github.com/RistoFlink/bookings/internal/helpers"
	"github.com/RistoFlink/bookings/internal/models"
	"github.com/RistoFlink/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig

var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Starting the server at: %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// what will be put into the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// set up an info logger and make it part of the app config
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	//set up an error logger and make it part of the app config
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	helpers.NewHelpers(&app)

	return nil
}
