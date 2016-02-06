package git

type GitInterface interface {
	FetchRepository(string) (string, error)
	addToCache(string, string)
	getFromCache(string) (string, bool)
}

type GitConfig struct {
	Priority int
	Path     string
}

type GitBackend struct {
	Repositories map[string]string
}

func (gb *GitBackend) addToCache(repo string, cloneURL string) {
	gb.Repositories[repo] = cloneURL
}

func (gb *GitBackend) getFromCache(repo string) (string, bool) {

	if gb.Repositories == nil {
		gb.Repositories = make(map[string]string)
	} else if cloneURL, ok := gb.Repositories[repo]; ok {
		return cloneURL, ok
	}
	return "", false
}

func GetRepository(gi GitInterface, repo string) (string, error) {

	cloneURL, ok := gi.getFromCache(repo)
	if !ok {
		var err error
		cloneURL, err = gi.FetchRepository(repo)
		if err != nil {
			return "", err
		}
		gi.addToCache(repo, cloneURL)
	}

	return cloneURL, nil
}
