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
	"bytes"
	"fmt"
	"io"
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
	mux.HandleFunc("/request/retrieve", retrieveRequests)
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

func retrieveRequests(w http.ResponseWriter, r *http.Request) {
	log.Print("Creating a request")
	client := &http.Client{}
	reqUrl := fmt.Sprintf("%s/api/ipwc/request_item/create/88858a471b44fbc4f141a8217e4bcbec/ritm_nv", snowenv)

	jsonData := `{
	"requested_for": "hdeshpande006",
		"variables_user": {
		"requested_for": "hdeshpande006"
	},
	"variables": {
		"departmentlos": "IFS",
			"primaryowner": "chintan.t.shah@au.pwc.com",
			"reqtype": " advreq",
			"secondaryowner": "6c23ed2ddb1cc380e61a384c7c96198b",
			"application_name": "4ecfe5ccdb96734009cd9044db96198f",
			"application_name_ci_number": "CI17511460"
			}
	}`

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer([]byte(jsonData)))
	req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)
	if err != nil {
		log.Fatalf("Error in initial request: %s", err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error in initial request: %s", err)
	}
	log.Println("Producing output")
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
	defer resp.Body.Close()
	w.Write([]byte(data))
}
