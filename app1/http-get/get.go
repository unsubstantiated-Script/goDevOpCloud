package http_get

import (
	"encoding/json"
	"fmt"
	"goDevOpCloud/utils"
	"io"
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

// Response This interface helps Words and Occurrence process through w/o hassle between types. Represents an intersection.
type Response interface {
	GetResponse() string
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

// DoGetRequest This API struct is in the init.go file
func (a API) DoGetRequest(requestURL string) (Response, error) {

	////Need to declare here else, the inline below only has acces to the var inside of that if scope.
	//var resp *http.Response

	resp, err := a.Client.Get(requestURL)

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
