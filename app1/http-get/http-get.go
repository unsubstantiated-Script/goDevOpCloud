package http_get

import (
	"flag"
	"fmt"
	"goDevOpCloud/utils"
	"net/url"
	"os"
)

func RollHTTPGet() {
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

	apiInstance := New(Options{
		Password: password,
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
	})

	res, err := apiInstance.DoGetRequest(parsedURL.String())
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
