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
			payload:   &azuredevops.PullRequestResource{},
			eventType: "git.pullrequest.created",
		},
		{
			payload:   azuredevops.PullRequestResource{},
			eventType: "git.pullrequest.merged",
		},
		{
			payload:   azuredevops.PullRequestResource{},
			eventType: "git.pullrequest.updated",
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
		t.Fatalf("ParsePayload: %+v, %v", got, err)
		if err != nil {
			t.Fatalf("ParsePayload: %v", err)
		}
		if want := test.payload; !reflect.DeepEqual(got, want) {
			t.Errorf("ParseWebHook(%#v, %#v) = %#v, want %#v", test.eventType, payload, got, want)
		}
	}
}
