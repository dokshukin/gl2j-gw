package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	BindIP       string `json:"bind_ip,omitempty" yaml:"bind_ip,omitempty"`
	BindPort     string `json:"bind_port,omitempty" yaml:"bind_port,omitempty"`
	Domain       string `json:"domain,omitempty" yaml:"domain,omitempty"`
	APIURI       string `json:"api_uri,omitempty" yaml:"api_uri,omitempty"`
	JiraSettings struct {
		URL      string `json:"url" yaml:"url"`
		User     string `json:"user" yaml:"user"`
		Password string `json:"password" yaml:"password"`
	} `json:"jira_settings" yaml:"jira_settings"`
	Projects map[string]project `json:"projects,omitempty" yaml:"projects,omitempty"`
}

type project struct {
	WorkFlow map[string]workflow `yaml:",inline"`
}

type workflow struct {
	Transitions struct {
		Created      int `json:"branch_created,omitempty" yaml:"branch_created,omitempty"`
		Merged       int `json:"merged,omitempty" yaml:"merged,omitempty"`
		MergeRequest int `json:"merge_request,omitempty" yaml:"merge_request,omitempty"`
	} `json:"transitions,omitempty" yaml:"transitions,omitempty"`
	ExternalLinks struct {
		Created      bool `json:"branch_created,omitempty" yaml:"branch_created,omitempty"`
		Deleted      bool `json:"branch_deleted,omitempty" yaml:"branch_deleted,omitempty"`
		Pipeline     bool `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
		Approved     bool `json:"approved,omitempty" yaml:"approved,omitempty"`
		MergeRequest bool `json:"merge_request,omitempty" yaml:"merge_request,omitempty"`
	} `json:"external_links,omitempty" yaml:"external_links,omitempty"`
}

func readConfig(configFile *string) (err error) {

	// read bytes from file
	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return
	}

	// parse YAML from []byte
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return
	}

	// exclude wrong URL settings fetched from config
	cfg.APIURI = strings.TrimSpace(cfg.APIURI)
	cfg.APIURI = strings.TrimRight(cfg.APIURI, "/")
	cfg.JiraSettings.URL = strings.TrimSpace(cfg.JiraSettings.URL)
	cfg.JiraSettings.URL = strings.TrimRight(cfg.JiraSettings.URL, "/")

	// Check mandatory parameters in config.
	// Exit, if they do not exist.
	if len(cfg.JiraSettings.URL) == 0 {
		err = errors.New("wrong Jira URL, please set up config")
		return
	}
	if len(cfg.JiraSettings.User) == 0 {
		err = errors.New("wrong Jira user, please set up config")
		return
	}
	if len(cfg.JiraSettings.Password) == 0 {
		err = errors.New("wrong Jira password, please set up config")
		return
	}
	return
}
