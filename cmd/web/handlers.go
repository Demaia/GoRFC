package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	supporting "gorfc.demaia.io/internal"
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

func Home(w http.ResponseWriter, r *http.Request) {
	log.Print("Homepage")
}

func CreateChange(w http.ResponseWriter, r *http.Request) {
	templateSysId := r.PathValue("tmpl")
	organisation := r.Header.Get("organisation")
	project := r.Header.Get("project")
	pipeline := r.Header.Get("pipeline")
	run := r.Header.Get("run")
	appci := r.Header.Get("appci")

	// definitionid := r.Header.Get("definitionid")
	// Start Create change
	client := &http.Client{}
	requrl := fmt.Sprintf("%s/api/sn_chg_rest/change/standard/%s", snowenv, templateSysId)
	dates := supporting.StartEnd()
	appIdentifier := supporting.RetrieveCiSysId(appci)
	details := fmt.Sprintf(`{
        "assignment_group": "f5ce7812db1a841084055ad6dc96197c",
        "u_coordinator_group": "f5ce7812db1a841084055ad6dc96197c",
        "assigned_to": "c6e5660e8754c6506e3462cbbbbb35b0",
        "u_change_manager": "c6e5660e8754c6506e3462cbbbbb35b0",
        "cmdb_ci": %s,
        "start_date": "%s",
        "requested_by_date": "%s",
        "end_date": "%s"
	}`, appIdentifier, dates["start"], dates["end"], dates["end"])

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

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
	displayname := supporting.RetrieveApprover("https://dev.azure.com/PwC-NL-APPS/", "Cloud%20Solutions%20Platform")
	supporting.AddWorknotes(organisation, project, pipeline, run, displayname, ChgCreate.Result.SysID.Value)
	supporting.MoveToImplement(ChgCreate.Result.SysID.Value)
	_, err = w.Write([]byte(ChgCreate.Result.SysID.Value))
	if err != nil {
		return
	}

}

func RetrieveChangeNo(w http.ResponseWriter, r *http.Request) {
}

func CloseChange(w http.ResponseWriter, r *http.Request) {
	closeNotes := fmt.Sprintf(`{
        "close_notes": "Change successful.",
        "u_close_code": "Change Successful"
    }`)
	chgsysid := r.PathValue("chgid")

	client := &http.Client{}
	reqURL := fmt.Sprintf("%s/api/sn_chg_rest/change/standard/%s", snowenv, chgsysid)
	req, err := http.NewRequest("PATCH", reqURL, bytes.NewBuffer([]byte(closeNotes)))
	if err != nil {
		log.Printf(err.Error())
	}
	req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf(err.Error())
	}
	defer resp.Body.Close()
}

func CancelChange(w http.ResponseWriter, r *http.Request) {
}

func RetrieveInc(w http.ResponseWriter, r *http.Request) {
	log.Print("Creating a request")
	client := &http.Client{}
	reqUrl := fmt.Sprintf("%s/api/now/table/incident", snowenv)

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

func RetrieveCi(w http.ResponseWriter, r *http.Request) {
	appci := r.PathValue("appci")

	client := &http.Client{}
	requrl := fmt.Sprintf("%s/cmdb_ci_service_discovered_list.do?JSONv2&sysparm_action=getRecords&sysparm_query=u_number=%s", snowenv, appci)
	req, err := http.NewRequest("GET", requrl, bytes.NewBuffer([]byte("")))

	req.SetBasicAuth(SnowServiceAccountName, SnowServiceAccountPassword)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if err != nil {
		log.Fatalf("Error in initial request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error in initial request: %s", err)
	}
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error in initial request: %s", err)
	}
	fmt.Println(string(data))
}
