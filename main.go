package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
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
		iter := jsoniter.Parse(jsoniter.ConfigFastest, rc, 1000)
		// read open bracket

		// while the array contains values
		for iter.ReadArray() {
			var gc GitCommit
			// decode an array value (GitCommit)

			obj := iter.Read()
			config := &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.StringToTimeHookFunc("2006-01-02T15:04:05Z"),
				Result:     &gc,
			}
			dec, err := mapstructure.NewDecoder(config)
			if err != nil {
				log.Fatal(err)
			}
			err = dec.Decode(obj)
			if err != nil {
				log.Fatal(err)
			}

			pw.Write([]byte(fmt.Sprintf("%v",
				struct {
					name  string
					email string
					time  string
				}{
					gc.Commit.Author.Name,
					gc.Commit.Author.Email,
					gc.Commit.Author.Date.Format(time.RFC3339),
				})))

		}

	}()

	return pr, nil
}

func write(trc io.ReadCloser) error {
	io.Copy(os.Stdout, trc)
	trc.Close()
	return nil
}
