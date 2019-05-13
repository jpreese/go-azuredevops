// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted for Azure Devops

package azuredevops_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mcdafydd/go-azuredevops/azuredevops"
)

func TestParsePayload(t *testing.T) {
	tests := []struct {
		payload   interface{}
		eventType string
	}{
		{
			payload:   &azuredevops.GitPullRequest{},
			eventType: "git.pullrequest.created",
		},
		{
			payload:   &azuredevops.GitPullRequest{},
			eventType: "git.pullrequest.merged",
		},
		{
			payload:   &azuredevops.GitPullRequest{},
			eventType: "git.pullrequest.updated",
		},
		{
			payload:   &azuredevops.GitPush{},
			eventType: "git.push",
		},
		{
			payload:   &azuredevops.WorkItem{},
			eventType: "workitem.commented",
		},
		{
			payload:   &azuredevops.WorkItemUpdate{},
			eventType: "workitem.updated",
		},
	}

	for _, test := range tests {
		event := new(azuredevops.Event)
		event.EventType = test.eventType
		payload, err := json.Marshal(test.payload)
		event.RawPayload = (json.RawMessage)(payload)
		if err != nil {
			t.Fatalf("Marshal(%#v): %v", test.payload, err)
		}
		got, err := event.ParsePayload()
		if err != nil {
			t.Fatalf("ParsePayload: %v", err)
		}
		if want := test.payload; !cmp.Equal(got, want) {
			diff := cmp.Diff(got, want)
			t.Errorf("ParsePayload error: %s", diff)
		}

	}
}
