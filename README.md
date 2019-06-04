# GitLab to Jira gateway [![Build Status](https://travis-ci.org/dokshukin/gl2j-gw.svg)][travis]
* issue transition from Gitlab to Jira
* creation of external links in jira to GitLab


## Documentation
[Quick start with docker](https://github.com/dokshukin/gl2j-gw/wiki/Quickstart)

[Detailed configuration](https://github.com/dokshukin/gl2j-gw/wiki/Configuration)

## Configure and run

### Download
For Linux (amd64):

    wget https://github.com/dokshukin/gl2j-gw/releases/download/v0.1/gl2j-gw_linux-amd64-v0.1.14 \
      -o /dev/null -O gl2j-gw && chmod +x gl2j-gw

See other distributions in [downloads](https://github.com/dokshukin/gl2j-gw/releases).

### Create config
Example (`config.yml`):
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
* docker containers
* WEB inerface (maybe)