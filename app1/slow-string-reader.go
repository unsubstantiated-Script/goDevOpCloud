package app1

import (
	"fmt"
	"io"
	"log"
)

// MySlowReader This sucker is gonna read an entire string char by char
type MySlowReader struct {
	contents string //Contents to be read
	pos      int    //Position in the string
}

func (m *MySlowReader) Read(p []byte) (n int, err error) {
	// Making sure the reader doesn't go any further than the length of the contents.
	if m.pos+1 <= len(m.contents) {
		//Copying just one single character in the string.
		n := copy(p, m.contents[m.pos:m.pos+1])
		m.pos++
		return n, nil
	}
	return 0, io.EOF
}

func SlowStringReader() {

	mySlowReaderInstance := &MySlowReader{
		contents: "stupid stuff",
	}

	out, err := io.ReadAll(mySlowReaderInstance)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("output: %s\n", out)
}
