package main

import (
	"fmt"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/debug/", debugSource())
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(mux)))
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile("\t(?:[A-Z]:|/)(.*\\.go):([0-9]+)")
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				stack = re.ReplaceAll(stack, []byte("\t<a href=\"/debug/$1?line=$2\">/$1:$2</a>"))
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, stack)
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func debugSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/debug") {
			path = strings.TrimPrefix(path, "/debug")
			data, err := os.ReadFile(path)
			if err != nil {
				http.Error(w, fmt.Sprintf("Problem reading file, %v", err), http.StatusInternalServerError)
				return
			}
			opts := []html.Option{html.Standalone(true), html.WithLineNumbers(true)}
			lineNumber := r.URL.Query().Get("line")
			if line, err := strconv.Atoi(lineNumber); err == nil {
				opts = append(opts, html.HighlightLines([][2]int{
					{line, line},
				}))
			}

			content := string(data)
			lexer := lexers.Analyse(content)
			tok, err := lexer.Tokenise(nil, content)
			if err != nil {
				http.Error(w, "Problem tokenizing file", http.StatusInternalServerError)
				return
			}
			f := html.New(opts...)
			s := styles.Get("dracula")
			if s == nil {
				s = styles.Fallback
			}
			w.WriteHeader(http.StatusOK)
			f.Format(w, s, tok)
		}
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
