package app1

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

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

	fmt.Printf("HTTP Status Code: %d\nBody: %s\n", resp.StatusCode, body)

}
