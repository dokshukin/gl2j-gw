# GitLab to Jira gateway [![Build Status](https://travis-ci.org/dokshukin/gl2j-gw.svg)][travis]
* issue transition from Gitlab to Jira
* creation of external links in jira to GitLab

## Configure and run

### Download
    wget 


### Create config
Config example
```
---

bind_ip: "0.0.0.0"
bind_port: "8080"
api_uri: "/api"

jira_settings:
  url: "https://your-domain.atlassian.net"
  user: "jira-bot@your-domain.com"
  password: "XXXXXXXXXXXXXXXXXXXX"

projects:
  DEV:
    Bug:
      transitions:
        branch_created: 21
        merged: 31
        merge_request: 41
      external_links:
        branch_created: True
        branch_deleted: True
        pipeline: False
        approved: True
        merge_request: True
    Feature:
      transitions:
        branch_created: 21
        merged: 31
        merge_request: 41
      external_links:
        branch_created: True
        branch_deleted: True
        pipeline: True
        approved: True
        merge_request: True
  OPS:
    Task:
      transitions:
        branch_created: 23
        merged: 33
        merge_request: 43
      external_links:
        branch_created: True
        branch_deleted: True
        pipeline: True
        approved: True
        merge_request: True
    INC:
      transitions:
        branch_created: 23
        merged: 33
        merge_request: 43
      external_links:
        branch_created: True
        branch_deleted: True
        pipeline: True
        approved: False
        merge_request: True
```

### Run
    ./gl2j-gw --config=/path/to/config.yml

## ToDo
* handle delete branch
* additional statuses for transition (f.e. `approved`, `tested`)
* branch creation from Jira


