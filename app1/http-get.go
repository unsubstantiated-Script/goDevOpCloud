package app1

import (
	"encoding/json"
	"flag"
	"fmt"
	"goDevOpCloud/utils"
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
	//args := os.Args
	//
	//if len(args) < 2 {
	//	fmt.Printf("Useage: ./http-get <url>\n")
	//	os.Exit(1)
	//}

	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)

	//Using flags to access CLI commands or endpoints
	flag.StringVar(&requestURL, "url", "", "url to access")
	flag.StringVar(&password, "password", "", "use a password to access API")

	flag.Parse()

	//Declaring and checking at the same time.
	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Validation error: URL is not valid %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	client := http.Client{}

	if password != "" {
		token, err := doLoginRequest(client, parsedURL.Scheme+"://"+parsedURL.Host+"/login", password)
		if requestErr, ok := err.(utils.RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		client.Transport = MyJWTTransport{
			transport: http.DefaultTransport,
			token:     token,
		}
	}

	res, err := doRequest(client, parsedURL.String())
	if err != nil {
		if requestErr, ok := err.(utils.RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No response \n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}

func doRequest(client http.Client, requestURL string) (Response, error) {

	////Need to declare here else, the inline below only has acces to the var inside of that if scope.
	//var resp *http.Response

	resp, err := client.Get(requestURL)

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

	if !json.Valid(body) {
		return nil, utils.RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, utils.RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Unmarshal error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, utils.RequestError{
				HTTPCode: resp.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words unmarshal error: %s", err),
			}
		}
		return words, nil
	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, utils.RequestError{
				HTTPCode: resp.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurrences unmarshal error: %s", err),
			}
		}

		return occurrence, nil
	}

	return nil, nil

}
