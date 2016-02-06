package bitbucket

import (
	"fmt"
	"net/http"
)

type RepositoriesService struct {
	client *Client
}

type Repository struct {
	Name    string          `json:"name"`
	Owner   RepositoryOwner `json:"owner"`
	Scm     string          `json:"scm"`
	Private bool            `json:"is_private"`
	Links   RepositoryLinks `json:"links"`
}

type RepositoryOwner struct {
	Username string `json:"username"`
}

type RepositoryLinks struct {
	CloneURL []RepositoryClone `json:"clone"`
}

type RepositoryClone struct {
	URL  string `json:"href"`
	Name string `string:"name"`
}

func (s *RepositoriesService) Get(owner, repo string) (*Repository, *http.Response, error) {

	u := fmt.Sprintf("repositories/%v/%v", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repository := new(Repository)
	resp, err := s.client.Do(req, repository)
	if err != nil {
		return nil, resp, err
	}

	return repository, resp, err
}
