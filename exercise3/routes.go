package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type server struct {
	router *http.ServeMux
	story  *Story
}

func newServer() *server {
	s := &server{}
	s.routes()
	return s
}

func (s *server) routes() {
	s.router = http.NewServeMux()
	s.router.Handle("/", s.pathLogger(s.handleStory()))
	s.router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleStory() http.HandlerFunc {
	tmpl, tmplErr := template.ParseFiles("templates/home.html")

	return func(w http.ResponseWriter, r *http.Request) {
		if tmplErr != nil {
			http.Error(w, "Encountered error", http.StatusInternalServerError)
		}
		parts := parsePath(r.URL.Path)
		if parts != nil {
			v, ok := (*s.story)[parts[0]]
			if ok {
				tmpl.Execute(w, v)
				return
			}
		}
		tmpl.Execute(w, (*s.story)["intro"])
	}
}

func (s *server) pathLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(">> Incoming request URL %s", r.URL.Path)
		h.ServeHTTP(w, r)
	}
}

func parsePath(path string) []string {
	if "/" == path {
		return nil
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	n := strings.Index(path, "/")
	if n == -1 {
		return []string{path}
	} else {
		return []string{path[:n], path[n+1:]}
	}
}
