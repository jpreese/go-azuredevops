package azuredevops_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/mcdafydd/go-azuredevops/azuredevops"
)

func TestPolicyEvaluationsService_List(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/o/p/_apis/policy/evaluations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"value": [
				{
				  "configuration": {
					"createdBy": {
						"displayName": "Normal Paulk",
						"url": "https://vssps.dev.azure.com/fabrikam/_apis/Identities/ac5aaba6-a66a-4e1d-b508-b060ec624fa9",
						"_links": {
							"avatar": {
								"href": "https://dev.azure.com/fabrikam/_apis/GraphProfile/MemberAvatars/aad.YmFjMGYyZDctNDA3ZC03OGRhLTlhMjUtNmJhZjUwMWFjY2U5"
							}
						},
						"id": "ac5aaba6-a66a-4e1d-b508-b060ec624fa9",
						"uniqueName": "dev@mailserver.com",
						"imageUrl": "https://dev.azure.com/fabrikam/_api/_common/identityImage?id=ac5aaba6-a66a-4e1d-b508-b060ec624fa9",
						"descriptor": "aad.YmFjMGYyZDctNDA3ZC03OGRhLTlhMjUtNmJhZjUwMWFjY2U5"
					},
					"createdDate": "2020-04-02T15:08:06.9682088Z",
					"isEnabled": true,
					"isBlocking": false,
					"isDeleted": false,
					"settings": {
					  "requiredReviewerIds": ["aaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"],
					  "minimumApproverCount": 1,
					  "creatorVoteCounts": true,
					  "scope": [
						{
						  "refName": "refs/heads/master",
						  "matchKind": "Exact",
						  "repositoryId": "aaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
						}
					  ],
					  "filenamePatterns": ["*"]
					},
					"isEnterpriseManaged": false,
					"_links": {
					  "self": {
						"href": "https://dev.azure.com/fabrikam/fa8735d5-fff3-464c-9012-1242593db299/_apis/policy/configurations/1873"
					  },
					  "policyType": {
						"href": "https://dev.azure.com/fabrikam/fa8735d5-fff3-464c-9012-1242593db299/_apis/policy/types/fd2167ab-b0be-447a-8ec8-39368250530e"
					  }
					},
					"revision": 13,
					"id": 1873,
					"url": "https://dev.azure.com/fabrikam/fa8735d5-fff3-464c-9012-1242593db299/_apis/policy/configurations/1873",
					"type": {
					  "id": "fd2167ab-b0be-447a-8ec8-39368250530e",
					  "url": "https://dev.azure.com/fabrikam/fa8735d5-fff3-464c-9012-1242593db299/_apis/policy/types/fd2167ab-b0be-447a-8ec8-39368250530e",
					  "displayName": "Required reviewers"
					}
				  },
				  "artifactId": "vstfs:///CodeReview/CodeReviewId/fa8735d5-fff3-464c-9012-1242593db299/17325",
				  "evaluationId": "2f582703-da53-43b0-ac96-81beb724f832",
				  "startedDate": "2020-04-03T12:54:11.6439378Z",
				  "status": "queued"
				}
			],
			"count": 1
		}`)
	})

	const artifactID = "vstfs:///CodeReview/CodeReviewId/fa8735d5-fff3-464c-9012-1242593db299/17325"

	got, _, err := c.PolicyEvaluations.List(context.Background(), "o", "p", artifactID, &azuredevops.PolicyEvaluationsListOptions{})
	if err != nil {
		t.Errorf("PolicyEvaluations.List returned error: %v", err)
	}

	if *got[0].ArtifactID != artifactID {
		t.Errorf("PolicyEvaluations.List ArtifactIDs don't match. got %v. want %v.", *got[0].ArtifactID, artifactID)
	}

	if *got[0].Configuration.IsBlocking != false {
		t.Errorf("PolicyEvaluations.List IsBlocking invalid return value. got %v. want %v.", *got[0].Configuration.IsBlocking, false)
	}
}
