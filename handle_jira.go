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
		Relationship: "GitLab",
		Application: jiraRemoteLinkApplication{
			Type: "com.gitlab",
			Name: "GitLab",
		},
		Object: jiraRemoteLinkObject{
			Title: action,
			URL:   gitLabURL,
			Icon: jiraRemoteLinkObjectIcon{
				URL:   cfg.Domain + gitLabLogo,
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
