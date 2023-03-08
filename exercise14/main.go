package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", PanicHandler(mux)))
}

func PanicHandler(inner http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var sb strings.Builder
				sb.Write(debug.Stack())
				log.Printf("Recovered from panic, error was: %s\n%s", r, sb.String())

				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("Something went wrong..."))
				if err != nil {
					log.Printf("Error writing reply: %s\n", err)
				}
			}
		}()

		inner.ServeHTTP(w, r)
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
