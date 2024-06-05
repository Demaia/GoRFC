/*
This application provides a simplified integration with ServiceNow

## Requirements Before this application can be succesfully deployed within an environment the following
environment variables will need to defined:

* SNOWSANAME - Name of the Service Now service account.
* SNOWSAPASS - Password of the Service Now service account
* APIKEY - Key for using the Service Now api
* APIPASS - Password for the Service Now api
* SNOWENV - Service now environment URL
* AUTORFC_PROXYAUTH - Authentication string for additional proxy

## Capabilities

Using this application you will be able to perform the following tasks:

* Create a change
* Advance the change to the implement stage.
* Close the change when successful.
* Close the change when unsuccesful (_Not yet implemented_).
*/
package main

import (
	"log"
	"net/http"
	"os"
)

var (
	SnowServiceAccountName     string            = os.Getenv("AUTORFC_SNOWSANAME")
	SnowServiceAccountPassword string            = os.Getenv("AUTORFC_SNOWSAPASS")
	snowenv                    string            = os.Getenv("AUTORFC_SNOWENV")
	template_sys_id            string            = "f7cfa23fdb39421052e652f3f396192f"
	headers                    map[string]string = map[string]string{
		"apikey":              os.Getenv("AUTORFC_APIKEY"),
		"apikeysecret":        os.Getenv("AUTORFC_APIPASS"),
		"Proxy-Authorization": os.Getenv("AUTORFC_PROXYAUTH"),
		"Accept":              "application/json",
		"Content-Type":        "application/json",
	}
)

func main() {
	log.Printf("Service account: 	%s", SnowServiceAccountName)
	log.Printf("SA password: 		%s", SnowServiceAccountPassword)
	log.Printf("Environment in use: %s", snowenv)
	log.Printf("Template in use: 	%s", template_sys_id)
	for k, v := range headers {
		log.Printf("%s : %s", k, v)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/change/create", createChange)
	mux.HandleFunc("/change/retrieve", retrieveChangeNo)
	mux.HandleFunc("/change/implement", implementChange)
	mux.HandleFunc("/change/close", closeChange)
	mux.HandleFunc("/change/cancel", cancelChange)
	mux.HandleFunc("/change/closectask", closeChangeCtask)
	log.Print("Starting the server at :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Print("Homepage")
}

func createChange(w http.ResponseWriter, r *http.Request) {
}

func retrieveChangeNo(w http.ResponseWriter, r *http.Request) {

}

func implementChange(w http.ResponseWriter, r *http.Request) {

}

func closeChange(w http.ResponseWriter, r *http.Request) {

}

func cancelChange(w http.ResponseWriter, r *http.Request) {

}
func closeChangeCtask(w http.ResponseWriter, r *http.Request) {

}
