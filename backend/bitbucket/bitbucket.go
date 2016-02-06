package bitbucket

import (
	"errors"
	"net/http"

	"github.com/0rax/go-redirect/backend/git"
	"github.com/mrjones/oauth"
)

type BitbucketConfig struct {
	git.GitConfig
	OAuthKey    string
	OAuthSecret string
	Org         string
}

type BitbucketRepositories struct {
	git.GitBackend
	Cfg    BitbucketConfig
	client *Client
}

func (bb *BitbucketRepositories) createClient() {

	tc := http.DefaultClient
	if bb.Cfg.OAuthKey != "" && bb.Cfg.OAuthSecret != "" {
		ts := oauth.NewConsumer(bb.Cfg.OAuthKey, bb.Cfg.OAuthSecret, oauth.ServiceProvider{})
		accessToken := &oauth.AccessToken{}
		tc, _ = ts.MakeHttpClient(accessToken)
	}

	bb.client = NewClient(tc)
	bb.Repositories = make(map[string]string)
}

func (bb *BitbucketRepositories) FetchRepository(repoName string) (string, error) {

	if bb.client == nil {
		bb.createClient()
	}

	if bb.client == nil {
		bb.createClient()
	}
	repo, _, err := bb.client.Repositories.Get(bb.Cfg.Org, repoName)
	if err != nil {
		return "", err
	}
	for _, url := range repo.Links.CloneURL {
		if url.Name == "https" {
			return url.URL, nil
		}
	}
	return "", errors.New("Repositoru has no https fetch url")
}
