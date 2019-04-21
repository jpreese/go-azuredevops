package azuredevops

import (
	"fmt"
)

// FavouritesService handles communication with the favourites methods on the API
// So far it looks like this is undocumented, so this could change
type FavouritesService struct {
	client *Client
}

// FavouritesResponse describes the favourites response
type FavouritesResponse struct {
	Count      int          `json:"count"`
	Favourites []*Favourite `json:"value"`
}

// Favourite describes what a favourite is
type Favourite struct {
	ID           *string `json:"id,omitempty"`
	ArtifactName *string `json:"artifactName,omitempty"`
	ArtifactType *string `json:"artifactType,omitempty"`
	ArtifactID   *string `json:"artifactId,omitempty"`
}

// List returns a list of the favourite items from for the user
func (s *FavouritesService) List() ([]*Favourite, int, error) {
	URL := fmt.Sprintf(
		"_apis/Favorite/Favorites?artifactType=%s",
		"Microsoft.TeamFoundation.Git.Repository", // @todo This needs fixing
	)

	u, _ := s.client.BaseURL.Parse(URL)
	request, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}
	var response FavouritesResponse
	_, err = s.client.Execute(request, &response)

	return response.Favourites, response.Count, err
}
