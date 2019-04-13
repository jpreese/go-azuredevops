// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted for Azure Devops

package azuredevops_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/mcdafydd/go-azuredevops/azuredevops"
)

func TestParseWebHook(t *testing.T) {
	payload := azuredevops.Event{}
	p, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal(%#v): %v", payload, err)
	}
	got, err := azuredevops.ParseWebHook(p)
	if err != nil {
		t.Fatalf("ParseWebHook: %v", err)
	}
	if want := &payload; !reflect.DeepEqual(got, want) {
		t.Errorf("ParseWebHook(%#v) = %#v, want %#v", p, got, want)
	}
}

func TestActivityID(t *testing.T) {
	id := "8970a780-244e-11e7-91ca-da3aabcb9793"

	req, err := http.NewRequest("POST", "http://localhost", nil)
	if err != nil {
		t.Fatalf("ActivityID: %v", err)
	}
	req.Header.Set("X-VSS-ActivityID", id)

	got := azuredevops.GetActivityID(req)
	if got != id {
		t.Errorf("ActivityID(%#v) = %q, want %q", req, got, id)
	}
}

func TestSubscriptionID(t *testing.T) {
	id := "6b9490e4-940d-4d16-8dae-d36580e7e2b4"

	req, err := http.NewRequest("POST", "http://localhost", nil)
	if err != nil {
		t.Fatalf("SubscriptionID: %v", err)
	}
	req.Header.Set("X-VSS-SubscriptionId", id)

	got := azuredevops.GetSubscriptionID(req)
	if got != id {
		t.Errorf("SubscriptionID(%#v) = %q, want %q", req, got, id)
	}
}
