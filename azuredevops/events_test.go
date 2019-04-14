// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted for Azure Devops

package azuredevops_test

import (
	"encoding/json"
	"reflect"
	"testing"

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
	}

	for _, test := range tests {
		event := new(azuredevops.Event)
		event.EventType = &test.eventType
		payload, err := json.Marshal(test.payload)
		event.RawPayload = (*json.RawMessage)(&payload)
		if err != nil {
			t.Fatalf("Marshal(%#v): %v", test.payload, err)
		}
		got, err := event.ParsePayload()
		if err != nil {
			t.Fatalf("ParsePayload: %v", err)
		}
		if want := test.payload; !reflect.DeepEqual(got, want) {
			t.Errorf("ParsePayload(%#v, %#v) = %#v, want %#v", test.eventType, test.payload, got, want)
		}
	}
}
