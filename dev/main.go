package main

import (
	"crypto/rand"
	"io"
	"log"
	"net/http"
)

type customReader struct {
	iter int
	size int
}

func (r *customReader) Read(b []byte) (int, error) {
	maxRead := r.size - r.iter
	if maxRead == 0 {
		return 0, io.EOF
	}
	if maxRead > len(b) {
		maxRead = len(b)
	}
	n, err := rand.Read(b[:maxRead])
	r.iter += maxRead
	return n, err
}

func (r *customReader) Close() error {
	// Uncomment this for even more races :((
	// r.iter = 0
	return nil
}

func (r *customReader) Reset() {
	r.iter = 0
}

func main() {
	http.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {})
	go func() {
		if err := http.ListenAndServe("localhost:8060", nil); err != nil {
			log.Panicf("err: %v", err)
		}
	}()

	client := http.DefaultClient
	r := &customReader{size: 900000}

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8060", r)
	if err != nil {
		log.Printf("%v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode)
	}
	resp.Body.Close()

	r.Reset()
}
