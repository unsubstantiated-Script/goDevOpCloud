package app1

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Page struct {
	Name string `json:"page"`
}

// Words This is setup to hit the '/words' endpoint
type Words struct {
	//These need to be exported because the json parse method here needs it.
	Input string   `json:"input"`
	Words []string `json:"words"`
}

// Occurrence This is setup to hit the '/occurrence' endpoint
type Occurrence struct {
	Words map[string]int `json:"words"`
}

func RollHTTPGet() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Useage: ./http-get <url>\n")
		os.Exit(1)
	}

	//Declaring and checking at the same time.
	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("URL is in invalid format: %s\n", err)
		os.Exit(1)
	}

	////Need to declare here else, the inline below only has acces to the var inside of that if scope.
	//var resp *http.Response

	resp, err := http.Get(args[1])

	if err != nil {
		log.Fatal(err)
		//We don't need these below as they're not user errors
		//fmt.Printf("URL is in invalid format: %s\n", err)
		//os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP Status Code: %d\nBody: %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("JSON: Parsed\nPage: %s\nWords: %v\n", page.Name, strings.Join(words.Words, ", "))
	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			log.Fatal(err)
		}

		if val, ok := occurrence.Words["word3"]; ok {
			fmt.Printf("Found word1: %d\n", val)
		}

		for word, occ := range occurrence.Words {
			fmt.Printf("%s: %d\n", word, occ)
		}
	default:
		fmt.Printf("Page not found\n")
	}

}
