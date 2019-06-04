package main

type gitLabObjectPush struct {
	ObjectKind        string                    `json:"object_kind"`
	TotalCommitsCount int                       `json:"total_commits_count"`
	Ref               string                    `json:"ref"`
	After             string                    `json:"after"`
	Before            string                    `json:"before"`
	EventType         string                    `json:"event_type,omitempty"`
	ProjectID         int                       `json:"project_id"`
	Username          string                    `json:"user_name,omitempty"`
	Project           gitLabSubObjectProject    `json:"project,omitempty"`
	Repository        gitLabSubObjectrepository `json:"repository,omitempty"`
}

type gitLabObjectMergeRequest struct {
	ObjectKind       string                      `json:"object_kind"`
	User             gitLabSubObjectUser         `json:"user,omitempty"`
	ObjectAttributes gitLabSubObjectAttributesMR `json:"object_attributes,omitempty"`
	Project          gitLabSubObjectProject      `json:"project,omitempty"`
	Repository       gitLabSubObjectrepository   `json:"repository,omitempty"`
}

type gitLabSubObjectProject struct {
	ID                int    `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	WebURL            string `json:"web_url,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	GitSSHURL         string `json:"git_ssh_url,omitempty"`
	GitHTTPURL        string `json:"git_http_url,omitempty"`
	Namespace         string `json:"namespace,omitempty"`
	VisibilityLevel   int    `json:"visibility_level,omitempty"`
	PathWithNamespace string `json:"path_with_namespace,omitempty"`
	DefaultBranch     string `json:"default_branch,omitempty"`
	CIConfigPath      string `json:"ci_config_path,omitempty"`
	Homepage          string `json:"homepage,omitempty"`
	URL               string `json:"url,omitempty"`
	SSHURL            string `json:"ssh_url,omitempty"`
	HTTPURL           string `json:"http_url,omitempty"`
}

type gitLabSubObjectrepository struct {
	Name            string `json:"name,omitempty"`
	URL             string `json:"url,omitempty"`
	Description     string `json:"description,omitempty"`
	Homepage        string `json:"homepage,omitempty"`
	GitHTTPURL      string `json:"git_http_url,omitempty"`
	GitSSHURL       string `json:"git_ssh_url,omitempty"`
	VisibilityLevel int    `json:"visibility_level,omitempty"`
}

type gitLabSubObjectAttributesMR struct {
	State        string `json:"state,omitempty"`
	Action       string `json:"action,omitempty"`
	MergeStatus  string `json:"merge_status,omitempty"`
	SourceBranch string `json:"source_branch,omitempty"`
	URL          string `json:"url,omitempty"`
}

type gitLabSubObjectUser struct {
	Name      string `json:"name,omitempty"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}
