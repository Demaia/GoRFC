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
* Close the change when unsuccessful (_Not yet implemented_).
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type DispVal struct {
	DisplayValue string `json:"display_value,omitempty"`
	Value        string `json:"value,omitempty"`
}
type ChangeCreated struct {
	Result Result `json:"result,omitempty"`
}
type Result struct {
	Number DispVal `json:"number"`
	SysID  DispVal `json:"sys_id"`
	Meta   struct {
		IgnoredFields []interface{} `json:"ignoredFields"`
	} `json:"__meta"`
}
type Approver struct {
	Value []struct {
		Steps []struct {
			ActualApprover struct {
				DisplayName string `json:"displayName"`
			} `json:"actualApprover,omitempty"`
			Status string `json:"status"`
			//Status string `json:"status"`
		} `json:"steps"`
		//Status string `json:"status"`
	} `json:"value"`
}

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
	ChgCreate ChangeCreated
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
	organisation := r.Header.Get("organisation")
	project := r.Header.Get("project")
	pipeline := r.Header.Get("pipeline")
	run := r.Header.Get("run")
	// definitionid := r.Header.Get("definitionid")
	// Start Create change
	client := &http.Client{}
	requrl := fmt.Sprintf("%s/api/sn_chg_rest/change/standard/%s", snowenv, template_sys_id)
	dates := startEnd()
	details := fmt.Sprintf(`{
        "assignment_group": "f5ce7812db1a841084055ad6dc96197c",
        "u_coordinator_group": "f5ce7812db1a841084055ad6dc96197c",
        "assigned_to": "c6e5660e8754c6506e3462cbbbbb35b0",
        "u_change_manager": "c6e5660e8754c6506e3462cbbbbb35b0",
        "cmdb_ci": "b3291246db9b14143c01cde40596199e",
        "start_date": "%s",
        "requested_by_date": "%s",
        "end_date": "%s"
	}`, dates["start"], dates["end"], dates["end"])

	req, err := http.NewRequest("POST", requrl, bytes.NewBuffer([]byte(details)))
	req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)

	if err != nil {
		log.Print("Error starting first Post")
		log.Print(err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error in executing request: %s", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading return values")
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &ChgCreate)

	if err != nil {
		log.Println("Error unmarshalling response body")
		log.Print(err)
	}
	log.Print(ChgCreate.Result.Number.Value)
	log.Print(ChgCreate.Result.SysID.Value)
	// Stop Create change

	// Start retrieving approver
	displayname := retrieveApprover("https://dev.azure.com/PwC-NL-APPS/", "Cloud%20Solutions%20Platform")
	// Stop Retrieving approver
	addWorknotes(organisation, project, pipeline, run, displayname, ChgCreate.Result.SysID.Value)

	// Start Move to implement
	requrl = fmt.Sprintf("%s/api/sn_chg_rest/change/standard/%s", snowenv, ChgCreate.Result.SysID.Value)
	state := map[string]string{"new": "1", "implement": "-1"}
	for _, v := range state {
		states := map[string]string{"state": v}
		payload, _ := json.Marshal(states)
		req, err = http.NewRequest("PATCH", requrl, bytes.NewBuffer(payload))
		if err != nil {
			log.Print("error starting Patch request")
			log.Print(err)
		}
		req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		resp, err = client.Do(req)
		resp.Body.Close()

		if err != nil {
			log.Print("error executing Patch request")
			log.Print(err)
		}
	}
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

func startEnd() map[string]string {
	startTime := time.Now().Format("2006-01-02 15:04:05")
	stopTime := time.Now().Add(time.Hour).Format("2006-01-02 15:04:05")

	dates := map[string]string{"start": startTime, "end": stopTime}
	return dates
}

func retrieveApprover(organisation string, project string) string {
	var approver Approver
	azdoauth := os.Getenv("AZDO_AUTH")

	azDoHeaders := map[string]string{"Authorization": azdoauth}
	client := &http.Client{}
	reqUrl := fmt.Sprintf("%s%s/_apis/pipelines/approvals?$expand=steps&api-version=7.0-preview", organisation, project)
	req, _ := http.NewRequest("GET", reqUrl, bytes.NewBuffer([]byte(reqUrl)))
	for k, v := range azDoHeaders {
		req.Header.Add(k, v)
	}

	resp, _ := client.Do(req)
	data, err := io.ReadAll(resp.Body)
	// log.Print(data)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(data, &approver)

	if err != nil {
		log.Print(err)
	}
	name := approver.Value[0].Steps[0].ActualApprover.DisplayName
	return name

}

func addWorknotes(organisation string, project string, pipeline string, run string, displayname string, chgsysid string) {

	client := &http.Client{}
	requrl := fmt.Sprintf("%s/api/sn_chg_rest/change/%s", snowenv, chgsysid)
	payload := fmt.Sprintf(`{"work_notes": "Organisation: %s \n Project: %s \n Pipeline: %s \n Run: %s \n Approved by: %s \n"}`, organisation, project, pipeline, run, displayname)
	req, err := http.NewRequest("PATCH", requrl, bytes.NewBuffer([]byte(payload)))

	if err != nil {
		log.Print("Error starting first Post")
		log.Print(err)
	}
	req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	_, err = client.Do(req)

	if err != nil {
		log.Print("Error starting first Post")
		log.Print(err)
	}
}
