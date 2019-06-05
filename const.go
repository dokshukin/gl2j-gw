package main

// constant statuses of messages in queue
const (
	// https://docs.atlassian.com/software/jira/docs/api/REST/8.2.0/
	///rest/api/2/project/{projectIdOrKey}/statuses
	// jiraRestBranchPrefix          = "/tree"
	gitLabBranchPrefix            = "/tree"
	jiraRestSuffixRemoteLink      = "/remotelink"
	jiraRestSuffixTransitions     = "/transitions"
	jiraRestBaseProjectURI        = "/rest/api/2/project"
	jiraRestBaseIssueURI          = "/rest/api/2/issue"
	jiraRemoteLinkCreatedID       = "1"
	jiraRemoteLinkApprovedID      = "2"
	jiraRemoteLinkTestsPassedID   = "3"
	jiraRemoteLinkMergedRequestID = "4"
	jiraObjectKindPush            = "push"
	jiraObjectKindPipeline        = "pipeline"
	jiraObjectKindBuild           = "build"
	jiraObjectKindMergeRequest    = "merge_request"
	jiraObjectKindMerged          = "merged"

	maxAttempts = 30

	gitLabLogo       = "/img/gitlab.png"
	gitLabZeroCommit = "0000000000000000000000000000000000000000"
)
