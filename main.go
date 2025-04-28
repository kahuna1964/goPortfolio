package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/kahuna1964/goPortfolio/internal/app"
	"github.com/kahuna1964/goPortfolio/internal/routes"
)

// command line switches
// -port port no.  If not provided, the flag below will assign 8080 as the default

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port), // sPrintf returns a value (does not print out to std i/o)
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("...Server is Ready and listening on port %d\n", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
