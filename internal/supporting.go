package supporting

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
type Appci struct {
	Records []struct {
		Name  string `json:"name"`
		SysID string `json:"sys_id"`
	} `json:"records"`
}

var (
	SnowServiceAccountName     string = os.Getenv("AUTORFC_SNOWSANAME")
	SnowServiceAccountPassword string = os.Getenv("AUTORFC_SNOWSAPASS")
	snowenv                    string = os.Getenv("AUTORFC_SNOWENV")
	// template_sys_id            string            = "f7cfa23fdb39421052e652f3f396192f"
	headers map[string]string = map[string]string{
		"apikey":              os.Getenv("AUTORFC_APIKEY"),
		"apikeysecret":        os.Getenv("AUTORFC_APIPASS"),
		"Proxy-Authorization": os.Getenv("AUTORFC_PROXYAUTH"),
		"Accept":              "application/json",
		"Content-Type":        "application/json",
	}
	ChgCreate ChangeCreated
)

func StartEnd() map[string]string {
	startTime := time.Now().Format("2006-01-02 15:04:05")
	stopTime := time.Now().Add(time.Hour).Format("2006-01-02 15:04:05")

	dates := map[string]string{"start": startTime, "end": stopTime}
	return dates
}

func RetrieveApprover(organisation string, project string) string {
	// var approver Approver
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

func AddWorknotes(organisation string, project string, pipeline string, run string, displayname string, chgsysid string) {

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

func MoveToImplement(chgid string) {
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

func RetrieveCiSysId(appci string) string {
	var AppciSysId Appci
	client := &http.Client{}
	reqUrl := fmt.Sprintf("%s/cmdb_ci_service_discovered_list.do?JSONv2&sysparm_action=getRecords&sysparm_query=u_number=%s", snowenv, appci)
	fmt.Printf("URL: %s", reqUrl)
	fmt.Printf("URL: %s", appci)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Print(err)
	}
	req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(data, &AppciSysId)
	if err != nil {
		log.Print(err)
	}
	sysId := AppciSysId.Records[0].SysID
	resp.Body.Close()

	return sysId
}
