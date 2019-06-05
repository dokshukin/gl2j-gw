package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Envelope is universal json object to be parsed
type Envelope struct {
	ObjectKind string `json:"object_kind"`
	Msg        interface{}
}

// main function handles all POST requests
func postHandler(w http.ResponseWriter, r *http.Request) {

	// check method (POST only)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		closeConnection(w, nil)
		return
	}

	// get request body
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		closeConnection(w, e)
		return
	}

	// We do not know initial structure of JSON reated by GitLab.
	// But there is one solid field "object_kind" in any GitLab request.
	// possible values:
	// - "push"  (including branch creation and removal)
	// - "pipeline" (typically provides CI info)
	// - "build" (part of pipeline)
	// - "merge_requset" (includes actions approved, merge, update, open)
	// - "note" (commetns and discusstions)
	//
	// So first step is detection of "object_kind"
	// Second step unmarshal object to appropriate STRUCT
	var msg json.RawMessage
	env := Envelope{
		Msg: &msg,
	}

	errDecoder := json.Unmarshal(body, &env)

	// debug print full request body
	log.Println(string(body))

	// you can find values of "jiraObjectKindPush", "jiraObjectKindMergeRequest" and others in const.go
	if errDecoder == nil {
		var err error
		switch env.ObjectKind {

		case jiraObjectKindPush:
			err = handleGitlabPush(body)
			log.Println(err)

		case jiraObjectKindBuild:

		case jiraObjectKindPipeline:

		case jiraObjectKindMergeRequest:
			err = handleGitlabMergeRequest(body)
			log.Println(err)

		default:
			// log.Println(string(body))

		}
		closeConnection(w, err)
	} else {
		closeConnection(w, errDecoder)
	}

}

func closeConnection(w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
