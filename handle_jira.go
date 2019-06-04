package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
func sendJiraBranchUpdate(req gitLabObjectPush, status string) (err error) {

	branchName := strings.Split(req.Ref, "/")[2] // example: "refs/heads/DEV-42-gitlab2jira-status-transitions"

	// fetch jira issue id from branch name
	jProj, jIssueID, err := parseJiraIssueIDFromBranchName(branchName)
	if err != nil {
		log.Println(err)
		return
	}

	// issue properties
	issueSettings, err := getJiraIssueSettings(jProj, jIssueID)
	if err != nil {
		return
	}

	// if project and workflow exist in config file
	if proj, ok := cfg.Projects[jProj]; ok {
		log.Printf("debug proj: %v\n", proj)
		if wf, ok2 := proj.WorkFlow[issueSettings.Fields.IssueType.Name]; ok2 {
			if status == "create" {

				// create jira remote link
				if wf.ExternalLinks.Created == true {
					obj := jiraRemoteLink{
						GlobalID: "system=" + req.Project.WebURL + jiraRestBranchPrefix + "/" + branchName + "&id=" + jiraRemoteLinkCreatedID,
						Application: jiraRemoteLinkApplication{
							Type: "com.gitlab.branch",
							Name: "New branch",
						},
						Object: jiraRemoteLinkObject{
							Title:   branchName,
							Summary: "branch was created by " + req.Username,
							URL:     req.Project.WebURL + jiraRestBranchPrefix + "/" + branchName,
							Icon: jiraRemoteLinkObjectIcon{
								URL:   gitLabLogo,
								Title: "gitlab logo",
							},
						},
					}

					payload, _ := json.Marshal(obj)
					// create request URL
					requestURL := cfg.JiraSettings.URL + jiraRestBaseIssueURI + "/" + jProj + "-" + jIssueID + jiraRestSuffixRemoteLink
					// send data
					go jiraSendRequest("POST", requestURL, payload)
				}

				// make transition
				if wf.Transitions.Created > 0 {
					// 0 is default setting during
					requestURL := cfg.JiraSettings.URL + jiraRestBaseIssueURI + "/" + jProj + "-" + jIssueID + jiraRestSuffixTransitions
					var obj jiraTransition
					obj.Transition.ID = strconv.Itoa(wf.Transitions.Created)
					payload, _ := json.Marshal(obj)
					go jiraSendRequest("POST", requestURL, payload)
				}
			}
			if status == "delete" {
				// TODO: make it work
				// delete jira remote link
				if wf.ExternalLinks.Deleted == true {
					type deleteRemoteLink struct {
						GlobalID string `json:"globalId"`
					}
					var obj deleteRemoteLink
					obj.GlobalID = "system=" + req.Project.WebURL + jiraRestBranchPrefix + "/" + branchName + "&id=" + jiraRemoteLinkCreatedID

					payload, _ := json.Marshal(obj)
					// create request URL
					requestURL := cfg.JiraSettings.URL + jiraRestBaseIssueURI + "/" + jProj + "-" + jIssueID + jiraRestSuffixRemoteLink
					// send data
					go jiraSendRequest("DELETE", requestURL, payload)
				}
			}

		} else {
			err = errors.New("issue " + issueSettings.Fields.IssueType.Name + " doesn't exist in config")
		}
	} else {
		err = errors.New("project " + jProj + " doesn't exist in config")
	}

	return
}
*/

func parseJiraIssueIDFromBranchName(branchName string) (jProjID string, jIssueID string, err error) {
	branchNamePrepared := strings.ReplaceAll(branchName, "_", "-")
	branchNameArray := strings.Split(branchNamePrepared, "-")

	if len(branchNameArray) < 2 {
		err = errors.New("error parsing issue ID from '" + branchName + "'")
	} else {
		jProjID = branchNameArray[0]
		jIssueID = branchNameArray[1]
	}

	return
}

func jiraSendRequest(requestMethod string, requestURL string, payload []byte) {
	// prepare http client
	jiraReq, _ := http.NewRequest(requestMethod, requestURL, bytes.NewBuffer(payload))
	jiraReq.Header.Set("Content-Type", "application/json")
	jiraReq.SetBasicAuth(cfg.JiraSettings.User, cfg.JiraSettings.Password)

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// multiple attempts to send data
	for i := 0; i < maxAttempts; i++ {

		resp, err := netClient.Do(jiraReq)
		if err == nil {
			log.Println(resp)
			break
		} else {
			log.Println(err)
			// increasing delays for request
			time.Sleep(time.Duration(i*5) * time.Second)
		}
	}
}

func getJiraIssueSettings(proj string, issueID string) (issueSettings jiraIssue, err error) {

	// get issue settings
	url := cfg.JiraSettings.URL + jiraRestBaseIssueURI + "/" + proj + "-" + issueID
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cfg.JiraSettings.User, cfg.JiraSettings.Password)
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &issueSettings)
	return
}

// function sends "POST" queries to create a new RemoteLink in Jira
func sendJiraRemoteLink(action string, project string, issueID string, gitLabURL string) (err error) {
	url := cfg.JiraSettings.URL + jiraRestBaseIssueURI + "/" + project + "-" +
		issueID + jiraRestSuffixRemoteLink

	// struct of POST request body to Jira
	var obj jiraRemoteLink

	// set default values
	obj = jiraRemoteLink{
		Application: jiraRemoteLinkApplication{
			Type: "com.gitlab",
			Name: action,
		},
		Object: jiraRemoteLinkObject{
			Title: action,
			URL:   gitLabURL,
			Icon: jiraRemoteLinkObjectIcon{
				URL:   gitLabLogo,
				Title: "gitlab logo",
			},
		},
	}

	switch action {

	case "approved":
		obj.GlobalID = "system=" + cfg.JiraSettings.URL + gitLabBranchPrefix +
			"/" + project + "-" + issueID + "&id=" + jiraRemoteLinkApprovedID
		obj.Object.Summary = "Merge request was approved"

	case "branch created":
		obj.GlobalID = "system=" + cfg.JiraSettings.URL + gitLabBranchPrefix +
			"/" + project + "-" + issueID + "&id=" + jiraRemoteLinkCreatedID
		// obj.Object.Summary = gitLabURL

	case "merge request":
		obj.GlobalID = "system=" + cfg.JiraSettings.URL + gitLabBranchPrefix +
			"/" + project + "-" + issueID + "&id=" + jiraRemoteLinkMergedRequestID
		obj.Object.Summary = "created"

	default:
		err = errors.New("unknown action " + action + " to create Jira exrenal link")
		return
	}

	payload, err := json.Marshal(obj)
	if err != nil {
		return
	}

	go jiraSendRequest("POST", url, payload)
	return
}

func sendJiraTransition(id int, project string, issueID string) (err error) {
	url := cfg.JiraSettings.URL + jiraRestBaseIssueURI + "/" + project + "-" +
		issueID + jiraRestSuffixTransitions
	var obj jiraTransition
	obj.Transition.ID = strconv.Itoa(id)
	payload, err := json.Marshal(obj)
	if err != nil {
		return
	}
	go jiraSendRequest("POST", url, payload)
	return
}
