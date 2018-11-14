package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	rc, err := read()
	if err != nil {
		log.Fatal("Could not read")
	}
	trc, err := transform(rc)
	if err != nil {
		log.Fatal("Could not transform")
	}
	err = write(trc)
	if err != nil {
		log.Fatal("Could not write")
	}
}

func read() (io.ReadCloser, error) {
	return os.Stdin, nil
}

func transform(rc io.ReadCloser) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		dec := json.NewDecoder(rc)
		// read open bracket
		t, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v\n", t, t)

		// while the array contains values
		for dec.More() {
			var gc GitCommit
			// decode an array value (GitCommit)
			err := dec.Decode(&gc)
			if err != nil {
				log.Fatal(err)
			}
			pw.Write([]byte(fmt.Sprintf("%v", gc)))
		}

		// read closing bracket
		t, err = dec.Token()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return pr, nil
}

func write(trc io.ReadCloser) error {
	io.Copy(os.Stdout, trc)
	trc.Close()
	return nil
}
