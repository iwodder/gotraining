package main

import (
	"flag"
	"fmt"
	"gotraining/exercise4/link"
	"log"
	"os"
)

func main() {
	filename := flag.String("filename", "html1.html", "The HTML file we need to search in for links.")
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal("Exiting....couldn't open file....", err)
	}

	links, err := link.Parse(f)
	if err != nil {
		log.Fatal("Error encountered when extracting links", err)
	}
	for _, l := range links {
		fmt.Println(l)
	}
}
