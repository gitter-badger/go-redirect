package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/0rax/go-redirect/backend"
	"github.com/0rax/go-redirect/backend/git"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const goGetRedirectTemplate = `<!DOCTYPE html>
<html>
  <head>
    <title>Repository</title>
    <meta name="go-import" content="{{.Root}} {{.VCS}} {{.RedirectRoot}}">
  </head>
  <body>
    Content: {{.Root}} {{.VCS}} {{.RedirectRoot}}
  </body>
</html>`

func getRepository(repository string, path string) string {

	if path != "" {
		b, ok := backend.RepoPath[path]
		if !ok {
			return ""
		}
		url, err := git.GetRepository(b, repository)
		if err == nil && url != "" {
			return url
		}
		return ""
	}

	for _, p := range backend.RepoPriority {
		b := backend.RepoBackend[p]
		url, err := git.GetRepository(b, repository)
		if err == nil && url != "" {
			return url
		}
	}
	return ""
}

func querryRepository(w http.ResponseWriter, r *http.Request) {

	var goGet bool
	if get, ok := r.URL.Query()["go-get"]; ok && get[0] == "1" {
		goGet = true
	} else {
		goGet = false
	}

	params := mux.Vars(r)
	repository := params["repository"]
	path := params["path"]
	cloneURL := getRepository(repository, path)

	if cloneURL != "" && goGet {
		w.WriteHeader(http.StatusFound)
		t, err := template.New("goGetRedirect").Parse(goGetRedirectTemplate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		err = t.Execute(w, map[string]string{
			"Root":         "local.host/" + repository,
			"VCS":          "git",
			"RedirectRoot": cloneURL,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if cloneURL != "" {
		w.Header().Set("Location", cloneURL)
		w.WriteHeader(http.StatusSeeOther)
		fmt.Fprintf(w, "<a href=\"%s\">See Other</a>.\n", cloneURL)
	} else {
		w.WriteHeader(http.StatusNotFound)
		if path != "" {
			fmt.Fprintf(w, "404 page not found: Repository '%s' not found in '%s'\n", repository, path)
		} else {
			fmt.Fprintf(w, "404 page not found: Repository '%s' not found\n", repository)
		}
	}
}

func runWebServer(listen string) {

	n := negroni.Classic()

	router := mux.NewRouter()
	router.HandleFunc("/{path:[a-z0-9-]+}/{repository}", querryRepository).Methods("GET")
	router.HandleFunc("/{repository}", querryRepository).Methods("GET")

	n.UseHandler(router)
	n.Run(listen)
}
