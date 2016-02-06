package github

import (
	"net/http"

	"github.com/0rax/go-redirect/backend/git"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubConfig struct {
	git.GitConfig
	Org        string
	OAuthToken string
}

type GithubRepositories struct {
	git.GitBackend
	Cfg    GithubConfig
	client *github.Client
}

func (gh *GithubRepositories) createClient() {

	tc := http.DefaultClient
	if gh.Cfg.OAuthToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gh.Cfg.OAuthToken},
		)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	gh.client = github.NewClient(tc)
	gh.Repositories = make(map[string]string)
}

func (gh *GithubRepositories) FetchRepository(repoName string) (string, error) {
	if gh.client == nil {
		gh.createClient()
	}
	repo, _, err := gh.client.Repositories.Get(gh.Cfg.Org, repoName)
	if err != nil {
		return "", err
	}
	return *repo.CloneURL, nil
}
