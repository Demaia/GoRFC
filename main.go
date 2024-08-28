package main

import (
	"log"
	"net/http"
	"os"

	handlers "gorfc.demaia.io/cmd/web"
)

var (
	SnowServiceAccountName     string = os.Getenv("AUTORFC_SNOWSANAME")
	SnowServiceAccountPassword string = os.Getenv("AUTORFC_SNOWSAPASS")
	snowenv                    string = os.Getenv("AUTORFC_SNOWENV")
)

func main() {
	log.Printf("Service account: 	%s", SnowServiceAccountName)
	log.Printf("SA password: 		%s", SnowServiceAccountPassword)
	log.Printf("Environment in use: %s", snowenv)

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handlers.Home)
	mux.HandleFunc("/change/create/{tmpl}", handlers.CreateChange)
	mux.HandleFunc("/change/retrieve", handlers.RetrieveChangeNo)
	mux.HandleFunc("/change/close/{chgid}", handlers.CloseChange)
	mux.HandleFunc("/change/cancel", handlers.CancelChange)
	mux.HandleFunc("/incident/retrieve", handlers.RetrieveInc)
	mux.HandleFunc("/appci/{appid}", handlers.RetrieveCi)

	log.Print("Starting the server at :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
