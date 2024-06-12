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
type DispValBool struct {
	DisplayValue string `json:"display_value,omitempty"`
	Value        bool   `json:"value,omitempty"`
}
type DispValInt struct {
	DisplayValue string  `json:"display_value,omitempty"`
	Value        float32 `json:"value,omitempty"`
}
type DispValInternal struct {
	DisplayValue string `json:"display_value,omitempty"`
	Value        string `json:"value,omitempty"`
	Internal     string `json:"display_value_internal,omitempty"`
}
type ChangeCreated struct {
	Result Result `json:"result,omitempty"`
}
type Result struct {
	Reason                           DispVal         `json:"reason"`
	UConflictStatus                  DispVal         `json:"u_conflict_status"`
	UReasonForFailure                DispVal         `json:"u_reason_for_failure"`
	Parent                           DispVal         `json:"parent"`
	UConflictLastRun                 DispVal         `json:"u_conflict_last_run"`
	UOtherTeam                       DispValBool     `json:"u_other_team"`
	UReasonForFailureSubCategory     DispVal         `json:"u_reason_for_failure_sub_category"`
	WatchList                        DispVal         `json:"watch_list"`
	ProposedChange                   DispVal         `json:"proposed_change"`
	UDynamicSurvey                   DispVal         `json:"u_dynamic_survey"`
	UReferences                      DispVal         `json:"u_references"`
	XMiomsAzpipelineStageAttempt     DispVal         `json:"x_mioms_azpipeline_stage_attempt"`
	UponReject                       DispVal         `json:"upon_reject"`
	XMiomsAzpipelineTemplateID       DispVal         `json:"x_mioms_azpipeline_template_id"`
	UComponentDowntimeOnly           DispValBool     `json:"u_component_downtime_only"`
	UTest                            DispValBool     `json:"u_test"`
	SysUpdatedOn                     DispValInternal `json:"sys_updated_on"`
	Type                             DispVal         `json:"type"`
	ApprovalHistory                  DispVal         `json:"approval_history"`
	UChangeTemplate                  DispVal         `json:"u_change_template"`
	UCabDecision                     DispVal         `json:"u_cab_decision"`
	Skills                           DispVal         `json:"skills"`
	TestPlan                         DispVal         `json:"test_plan"`
	Number                           DispVal         `json:"number"`
	CabDelegate                      DispVal         `json:"cab_delegate"`
	UCiDescription                   DispVal         `json:"u_ci_description"`
	RequestedByDate                  DispValInternal `json:"requested_by_date"`
	CiClass                          DispVal         `json:"ci_class"`
	State                            DispValInt      `json:"state"`
	SysCreatedBy                     DispVal         `json:"sys_created_by"`
	Knowledge                        DispValBool     `json:"knowledge"`
	Order                            DispVal         `json:"order"`
	Phase                            DispVal         `json:"phase"`
	UPirRequired                     DispValBool     `json:"u_pir_required"`
	CmdbCi                           DispVal         `json:"cmdb_ci"`
	CmdbCiBusinessApp                DispVal         `json:"cmdb_ci_business_app"`
	Impact                           DispValInt      `json:"impact"`
	Contract                         DispVal         `json:"contract"`
	XMiomsAzpipelineBuildID          DispVal         `json:"x_mioms_azpipeline_build_id"`
	Active                           DispValBool     `json:"active"`
	WorkNotesList                    DispVal         `json:"work_notes_list"`
	UCabApprovers                    DispVal         `json:"u_cab_approvers"`
	Priority                         DispValInt      `json:"priority"`
	SysDomainPath                    DispVal         `json:"sys_domain_path"`
	ProductionSystem                 DispValBool     `json:"production_system"`
	CabRecommendation                DispVal         `json:"cab_recommendation"`
	RejectionGoto                    DispVal         `json:"rejection_goto"`
	ReviewDate                       DispValInternal `json:"review_date"`
	RequestedBy                      DispVal         `json:"requested_by"`
	GroupList                        DispVal         `json:"group_list"`
	BusinessDuration                 DispVal         `json:"business_duration"`
	UDateTaken                       DispValInternal `json:"u_date_taken"`
	ChangePlan                       DispVal         `json:"change_plan"`
	ApprovalSet                      DispValInternal `json:"approval_set"`
	ImplementationPlan               DispVal         `json:"implementation_plan"`
	UAffectedCiCount                 DispValInt      `json:"u_affected_ci_count"`
	UChangeManagerGroup              DispVal         `json:"u_change_manager_group"`
	UniversalRequest                 DispVal         `json:"universal_request"`
	EndDate                          DispValInternal `json:"end_date"`
	ShortDescription                 DispVal         `json:"short_description"`
	UCoordinatorGroup                DispVal         `json:"u_coordinator_group"`
	CorrelationDisplay               DispVal         `json:"correlation_display"`
	WorkStart                        DispValInternal `json:"work_start"`
	OutsideMaintenanceSchedule       DispValBool     `json:"outside_maintenance_schedule"`
	UDetailedReasonForFailure        DispVal         `json:"u_detailed_reason_for_failure"`
	AdditionalAssigneeList           DispVal         `json:"additional_assignee_list"`
	UActualCurrency                  DispVal         `json:"u_actual_currency"`
	StdChangeProducerVersion         DispVal         `json:"std_change_producer_version"`
	ServiceOffering                  DispVal         `json:"service_offering"`
	SysClassName                     DispVal         `json:"sys_class_name"`
	FollowUp                         DispValInternal `json:"follow_up"`
	ClosedBy                         DispVal         `json:"closed_by"`
	UPostChangeReviewRequired        DispValBool     `json:"u_post_change_review_required"`
	UParentImplementation            DispVal         `json:"u_parent_implementation"`
	URiskAndImpactAnalysis           DispVal         `json:"u_risk_and_impact_analysis"`
	ReviewStatus                     DispVal         `json:"review_status"`
	ReassignmentCount                DispValInt      `json:"reassignment_count"`
	UOpenedFor                       DispVal         `json:"u_opened_for"`
	URegion                          DispVal         `json:"u_region"`
	StartDate                        DispValInternal `json:"start_date"`
	AssignedTo                       DispVal         `json:"assigned_to"`
	Variables                        DispVal         `json:"variables"`
	SLADue                           DispValInternal `json:"sla_due"`
	UCoordinatorChangedBySystem      DispValBool     `json:"u_coordinator_changed_by_system"`
	CommentsAndWorkNotes             DispVal         `json:"comments_and_work_notes"`
	UOutcome                         DispVal         `json:"u_outcome"`
	UCoordinatorApproval             DispValBool     `json:"u_coordinator_approval"`
	UChangeCatalogue                 DispVal         `json:"u_change_catalogue"`
	Escalation                       DispVal         `json:"escalation"`
	UponApproval                     DispVal         `json:"upon_approval"`
	CorrelationID                    DispVal         `json:"correlation_id"`
	UCustomer                        DispVal         `json:"u_customer"`
	MadeSLA                          DispValBool     `json:"made_sla"`
	BackoutPlan                      DispVal         `json:"backout_plan"`
	UOpenedGroup                     DispVal         `json:"u_opened_group"`
	UActualEnd                       DispValInternal `json:"u_actual_end"`
	UTrackAg                         DispVal         `json:"u_track_ag"`
	UMissingCi                       DispValBool     `json:"u_missing_ci"`
	ConflictStatus                   DispVal         `json:"conflict_status"`
	TaskEffectiveNumber              DispVal         `json:"task_effective_number"`
	UAutomatedChangeNotificationsOff DispValBool     `json:"u_automated_change_notifications_off"`
	UPlannedDowntimeEnd              DispValInternal `json:"u_planned_downtime_end"`
	SysUpdatedBy                     DispVal         `json:"sys_updated_by"`
	USample                          DispValBool     `json:"u_sample"`
	UserInput                        DispVal         `json:"user_input"`
	OpenedBy                         DispVal         `json:"opened_by"`
	UUnauthorizedChange              DispValBool     `json:"u_unauthorized_change"`
	UQualityAssuranceCheckRequired   DispValBool     `json:"u_quality_assurance_check_required"`
	SysCreatedOn                     DispValInternal `json:"sys_created_on"`
	OnHoldTask                       DispVal         `json:"on_hold_task"`
	SysDomain                        DispVal         `json:"sys_domain"`
	UChangeApprovers                 DispVal         `json:"u_change_approvers"`
	UActualStart                     DispValInternal `json:"u_actual_start"`
	RouteReason                      DispVal         `json:"route_reason"`
	UIPAddress                       DispVal         `json:"u_ip_address"`
	ClosedAt                         DispValInternal `json:"closed_at"`
	UManuallyAddedAffectedCi         DispVal         `json:"u_manually_added_affected_ci"`
	UChangeManager                   DispVal         `json:"u_change_manager"`
	ReviewComments                   DispVal         `json:"review_comments"`
	BusinessService                  DispVal         `json:"business_service"`
	UOnHoldReason                    DispVal         `json:"u_on_hold_reason"`
	TimeWorked                       DispVal         `json:"time_worked"`
	ChgModel                         DispVal         `json:"chg_model"`
	ExpectedStart                    DispValInternal `json:"expected_start"`
	URepresentative                  DispVal         `json:"u_representative"`
	OpenedAt                         DispValInternal `json:"opened_at"`
	WorkEnd                          DispValInternal `json:"work_end"`
	PhaseState                       DispVal         `json:"phase_state"`
	XMiomsAzpipelineStageID          DispVal         `json:"x_mioms_azpipeline_stage_id"`
	UCreateModifyCi                  DispValBool     `json:"u_create_modify_ci"`
	UPlannedDowntimeStart            DispValInternal `json:"u_planned_downtime_start"`
	CabDate                          DispValInternal `json:"cab_date"`
	WorkNotes                        DispVal         `json:"work_notes"`
	UCloseChange                     DispVal         `json:"u_close_change"`
	CloseCode                        DispValBool     `json:"close_code"`
	AssignmentGroup                  DispVal         `json:"assignment_group"`
	XMiomsAzpipelinePipelineMetadata DispVal         `json:"x_mioms_azpipeline_pipeline_metadata"`
	UMissingCiName                   DispVal         `json:"u_missing_ci_name"`
	Description                      DispVal         `json:"description"`
	OnHoldReason                     DispVal         `json:"on_hold_reason"`
	UTakenBy                         DispVal         `json:"u_taken_by"`
	CalendarDuration                 DispVal         `json:"calendar_duration"`
	URetireCi                        DispValBool     `json:"u_retire_ci"`
	UUrgencyDriverReason             DispVal         `json:"u_urgency_driver_reason"`
	CloseNotes                       DispVal         `json:"close_notes"`
	UAutomatedCiCount                DispVal         `json:"u_automated_ci_count"`
	SysID                            DispVal         `json:"sys_id"`
	ContactType                      DispVal         `json:"contact_type"`
	CabRequired                      DispValBool     `json:"cab_required"`
	Urgency                          DispValInt      `json:"urgency"`
	Scope                            DispValInt      `json:"scope"`
	Company                          DispVal         `json:"company"`
	Justification                    DispVal         `json:"justification"`
	UStatus                          DispVal         `json:"u_status"`
	UChangeReason                    DispVal         `json:"u_change_reason"`
	ActivityDue                      DispValInternal `json:"activity_due"`
	Comments                         DispVal         `json:"comments"`
	UEnvironment                     DispVal         `json:"u_environment"`
	Approval                         DispVal         `json:"approval"`
	DueDate                          DispVal         `json:"due_date"`
	SysModCount                      DispValInt      `json:"sys_mod_count"`
	OnHold                           DispValBool     `json:"on_hold"`
	UReviewer                        DispVal         `json:"u_reviewer"`
	SysTags                          DispVal         `json:"sys_tags"`
	UPreventativeActions             DispVal         `json:"u_preventative_actions"`
	ConflictLastRun                  DispValInternal `json:"conflict_last_run"`
	CabDateTime                      DispValInternal `json:"cab_date_time"`
	RiskValue                        DispVal         `json:"risk_value"`
	URootCause                       DispVal         `json:"u_root_cause"`
	Unauthorized                     DispValBool     `json:"unauthorized"`
	UTrackAssignee                   DispVal         `json:"u_track_assignee"`
	UReasonForFailureCategory        DispVal         `json:"u_reason_for_failure_category"`
	UCloseCode                       DispVal         `json:"u_close_code"`
	UChangeAction                    DispVal         `json:"u_change_action"`
	Risk                             DispValInt      `json:"risk"`
	Location                         DispVal         `json:"location"`
	ULessonsLearnt                   DispVal         `json:"u_lessons_learnt"`
	UConfigurationItem               DispVal         `json:"u_configuration_item"`
	Category                         DispVal         `json:"category"`
	UCiManagingTerritory             DispVal         `json:"u_ci_managing_territory"`
	RiskImpactAnalysis               DispVal         `json:"risk_impact_analysis"`
	Meta                             struct {
		IgnoredFields []interface{} `json:"ignoredFields"`
	} `json:"__meta"`
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
		log.Print(err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Error in executing request: %s", err)
	}

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
	//fmt.Println(string(data))
	log.Print(ChgCreate.Result.Number.Value)
	log.Print(ChgCreate.Result.SysID.Value)

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
