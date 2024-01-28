package app1

import (
	"encoding/json"
	"fmt"
	"io"
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

// GetResponse method bound to above struct. Helps with our Response interface below
func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

// Occurrence This is setup to hit the '/occurrence' endpoint
type Occurrence struct {
	Words map[string]int `json:"words"`
}

// GetResponse method bound to above struct. Helps with our Response interface below
func (o Occurrence) GetResponse() string {
	//Transforming map into a slice
	var out []string

	for word, occur := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occur))
	}

	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

// Response This interface helps Words and Occurrence process through w/o hassle between types. Represents an intersection.
type Response interface {
	GetResponse() string
}

func RollHTTPGet() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Useage: ./http-get <url>\n")
		os.Exit(1)
	}

	res, err := doRequest(args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No response \n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}

func doRequest(requestURL string) (Response, error) {
	//Declaring and checking at the same time.
	if _, err := url.ParseRequestURI(requestURL); err != nil {
		return nil, fmt.Errorf("validation error: URL is not valid %s", err)
	}

	////Need to declare here else, the inline below only has acces to the var inside of that if scope.
	//var resp *http.Response

	resp, err := http.Get(requestURL)

	if err != nil {
		return nil, fmt.Errorf("HTTP Get error: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid output (HTTP Code %d): %s\n", resp.StatusCode, string(body))
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %s", err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %s", err)
		}
		return words, nil
	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %s", err)
		}

		return occurrence, nil
	}

	return nil, nil

}
