package main

type jiraRemoteLink struct {
	GlobalID     string                    `json:"globalId"`
	Relationship string                    `json:"relationship"`
	Application  jiraRemoteLinkApplication `json:"application"`
	Object       jiraRemoteLinkObject      `json:"object"`
}

type jiraRemoteLinkApplication struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type jiraRemoteLinkObject struct {
	URL     string                   `json:"url"`
	Title   string                   `json:"title"`
	Summary string                   `json:"summary"`
	Icon    jiraRemoteLinkObjectIcon `json:"icon"`
}

type jiraRemoteLinkObjectIcon struct {
	URL   string `json:"url16x16"`
	Title string `json:"title"`
}

type jiraIssue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		StatusCategoryChangeDate string `json:"statuscategorychangedate"`
		IssueType                struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Description string `json:"description"`
			Name        string `json:"name"`
			SubTask     bool   `json:"subtask"`
		} `json:"issuetype"`
		Project struct {
			ID             string `json:"id"`
			Self           string `json:"self"`
			Name           string `json:"name"`
			ProjectTypeKey string `json:"projectTypeKey"`
		} `json:"project"`
		Assignee struct {
			Self         string `json:"self"`
			Name         string `json:"name"`
			Key          string `json:"key"`
			AccountID    string `json:"accountId"`
			EmailAddress string `json:"emailAddress"`
			DisplayName  string `json:"displayName"`
			Active       bool   `json:"active"`
			TimeZone     string `json:"timeZone"`
			AccountType  string `json:"accountType"`
		} `json:"assignee"`
		Updated string `json:"updated"`
		Status  struct {
			ID          string `json:"id"`
			Self        string `json:"self"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"status"`
	} `json:"fields"`
}

type jiraTransition struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
}
