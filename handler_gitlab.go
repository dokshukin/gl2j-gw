package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

func handleGitlabPush(body []byte) (err error) {

	// structure of gitlab requests ("object_kind":"push")
	var req gitLabObjectPush

	if err = json.Unmarshal(body, &req); err != nil {
		return
	}
	branchName := strings.Split(req.Ref, "/")[2] // example: "refs/heads/DEV-42-gitlab2jira-status-transitions"

	// fetch jira issue id from branch name
	jProj, jIssueID, err := parseJiraIssueIDFromBranchName(branchName)
	if err != nil {
		return
	}
	// fetch jira isse settings
	issueSettings, err := getJiraIssueSettings(jProj, jIssueID)
	if err != nil {
		return
	}

	// check project in config
	if proj, ok := cfg.Projects[jProj]; ok {
		// check workflow in project in config
		if wf, ok2 := proj.WorkFlow[issueSettings.Fields.IssueType.Name]; ok2 {
			// branch creation is GitLab "push" request with sample fields:
			// { "before":"0000000000000000000000000000000000000000", "total_commits_count":0 }
			if req.Before == gitLabZeroCommit && req.TotalCommitsCount == 0 {
				// if config for the workflow requires creation of Remote Links in Jira
				if wf.ExternalLinks.Created == true {
					// new branch
					err = sendJiraRemoteLink("branch created", jProj, jIssueID, req.Project.WebURL+gitLabBranchPrefix+"/"+branchName)
					if err != nil {
						log.Println(err)
					}
				}
				// if config for the workflow requires transition in Jira
				if wf.Transitions.Created > 0 {
					err = sendJiraTransition(wf.Transitions.Created, jProj, jIssueID)
					if err != nil {
						log.Println(err)
					}
				}
			} else if req.After == gitLabZeroCommit {
				// err = sendJiraBranchUpdate(req, "delete")
				// log.Println(err)
			}
		}
	}
	return
}

func handleGitlabMergeRequest(body []byte) (err error) {

	// structure of gitlab requests ("object_kind":"merge_request")
	var req gitLabObjectMergeRequest
	if err = json.Unmarshal(body, &req); err != nil {
		return
	}

	// example { "object_attributes": {"source_branch": "DEV-42-hello-world"} }
	branchName := req.ObjectAttributes.SourceBranch

	// fetch jira issue id from branch name
	jProj, jIssueID, err := parseJiraIssueIDFromBranchName(branchName) // f.e.: ["DEV", "42", nil]
	if err != nil {
		return
	}

	// issue properties (makes request to Jira)
	issueSettings, err := getJiraIssueSettings(jProj, jIssueID)
	if err != nil {
		return
	}

	// check project in config
	if proj, ok := cfg.Projects[jProj]; ok {

		// check workflow in project in config
		if wf, ok2 := proj.WorkFlow[issueSettings.Fields.IssueType.Name]; ok2 {

			switch req.ObjectAttributes.Action {

			case "approved":
				if wf.ExternalLinks.Approved == true {
					err = sendJiraRemoteLink(req.ObjectAttributes.Action, jProj, jIssueID, req.ObjectAttributes.URL)
				}

			case "update":

			case "merge":
				// real status - branch was merged
				// Possible bug if Jira really supports transitions with ID == 0...
				// I haven't seen such IDs.
				// By default during initialization all int values are 0.
				if wf.Transitions.Merged > 0 {
					err = sendJiraTransition(wf.Transitions.Merged, jProj, jIssueID)
				}

			case "open":
				// merge requests was created
				if wf.Transitions.MergeRequest > 0 {
					err = sendJiraTransition(wf.Transitions.MergeRequest, jProj, jIssueID)
				}
				if wf.ExternalLinks.MergeRequest == true {
					err = sendJiraRemoteLink("merge request", jProj, jIssueID, req.ObjectAttributes.URL)
				}
			default:
				err = errors.New("Unknown action " + req.ObjectAttributes.Action + " in merge request structure")

			}
		} else {
			err = errors.New("no such issueType " + issueSettings.Fields.IssueType.Name + " in config file")
		}
	} else {
		err = errors.New("no such project " + jProj + " in config file")
	}
	return
}
