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

func moveToImplement(chgid string) {
	client := &http.Client{}
	requrl := fmt.Sprintf("%s/api/sn_chg_rest/change/standard/%s", snowenv, chgid)
	state := map[string]string{"new": "1", "implement": "-1"}
	for _, v := range state {
		states := map[string]string{"state": v}
		payload, _ := json.Marshal(states)
		req, err := http.NewRequest("PATCH", requrl, bytes.NewBuffer(payload))
		if err != nil {
			log.Print("error starting Patch request")
			log.Print(err)
		}
		req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)
		for k, v := range headers {
			req.Header.Add(k, v)
		}
		resp, err := client.Do(req)
		resp.Body.Close()

		if err != nil {
			log.Print("error executing Patch request")
			log.Print(err)
		}
	}
}
