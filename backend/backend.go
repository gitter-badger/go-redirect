package backend

import (
	"fmt"
	"sort"

	"github.com/0rax/go-redirect/backend/bitbucket"
	. "github.com/0rax/go-redirect/backend/git"
	"github.com/0rax/go-redirect/backend/github"
)

type Config struct {
	Github    []github.GithubConfig
	Bitbucket []bitbucket.BitbucketConfig
}

var (
	RepoBackend  map[int]GitInterface
	RepoPath     map[string]GitInterface
	RepoPriority []int
)

func addBackend(b string, cfg GitConfig, backend interface{}) {
	priority := cfg.Priority
	if priority <= 0 {
		priority = 100
	}
	for {
		if _, ok := RepoBackend[priority]; ok {
			priority++
			continue
		}
		break
	}
	RepoBackend[priority] = backend.(GitInterface)
	if cfg.Path != "" {
		RepoPath[cfg.Path] = backend.(GitInterface)
		fmt.Printf("[go-redirect] Added Valid Backend of type '%s' with priority: %d and path: %s\n", b, priority, cfg.Path)
	} else {
		fmt.Printf("[go-redirect] Added Valid Backend of type '%s' with priority: %d\n", b, priority)
	}

}

func Configure(backendConfig *Config) {

	RepoBackend = make(map[int]GitInterface)
	RepoPath = make(map[string]GitInterface)

	for _, bConfig := range backendConfig.Github {
		backend := github.GithubRepositories{
			Cfg: bConfig,
		}
		addBackend("github", bConfig.GitConfig, &backend)
	}
	for _, bConfig := range backendConfig.Bitbucket {
		backend := bitbucket.BitbucketRepositories{
			Cfg: bConfig,
		}
		addBackend("bitbucket", bConfig.GitConfig, &backend)
	}

	for k := range RepoBackend {
		RepoPriority = append(RepoPriority, k)
	}
	sort.Ints(RepoPriority)
}
