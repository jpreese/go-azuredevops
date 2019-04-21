package azuredevops

import (
	"fmt"
	"time"
)

// TeamsService handles communication with the teams methods on the API
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/core/teams/get%20all%20teams
type TeamsService struct {
	client *Client
}

// Project Describes a project
type Project struct {
	ID             *string    `json:"id,omitempty"`
	Name           *string    `json:"name,omitempty"`
	Description    *string    `json:"description,omitempty"`
	URL            *string    `json:"url,omitempty"`
	State          *string    `json:"state,omitempty"`
	Revision       *int       `json:"revision,omitempty"`
	Visibility     *string    `json:"visibility,omitempty"`
	LastUpdateTime *time.Time `json:"lastUpdateTime,omitempty"`
}

// Team describes what a team looks like
type Team struct {
	ID          string `url:"id,omitempty"`
	Name        string `url:"name,omitempty"`
	URL         string `url:"url,omitempty"`
	Description string `url:"description,omitempty"`
}

// TeamsListOptions describes what the request to the API should look like
type TeamsListOptions struct {
	Mine bool `url:"$mine,omitempty"`
	Top  int  `url:"$top,omitempty"`
	Skip int  `url:"$skip,omitempty"`
}

// TeamsListResponse Requests that may return multiple entities use this format
type TeamsListResponse struct {
	Count int     `json:"count,omitempty"`
	Teams []*Team `json:"value,omitempty"`
}

// TeamProjectCollectionReference Reference object for a TeamProjectCollection.
type TeamProjectCollectionReference struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

// TeamProjectReference Represents a shallow reference to a TeamProject.
type TeamProjectReference struct {
	Abbreviation        *string `json:"abbreviation,omitempty"`
	DefaultTeamImageURL *string `json:"defaultTeamImageUrl,omitempty"`
	Description         *string `json:"description,omitempty"`
	ID                  *string `json:"id,omitempty"`
	Name                *string `json:"name,omitempty"`
	Revision            *string `json:"revision,omitempty"`
	State               *string `json:"state,omitempty"`
	URL                 *string `json:"url,omitempty"`
	Visibility          *string `json:"visibility,omitempty"`
}

// List returns list of the teams
// https://docs.microsoft.com/en-us/rest/api/azure/devops/core/teams/get%20teams
// GET https://dev.azure.com/{organization}/_apis/projects/{projectId}/teams?api-version=5.1-preview.2
func (s *TeamsService) List(opts *TeamsListOptions) ([]*Team, int, error) {
	URL := fmt.Sprintf("_apis/teams?api-version=%s", APIVersion)
	URL, err := addOptions(URL, opts)

	u, _ := s.client.BaseURL.Parse(URL)
	request, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}
	var response TeamsListResponse
	_, err = s.client.Execute(request, &response)

	return response.Teams, response.Count, err
}
