package httpclient

import (
	"io/ioutil"
	"log"
	"net/http"
)

type addHeaderTransport struct {
	T http.RoundTripper
}

func (adt *addHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "go")
	return adt.T.RoundTrip(req)
}

func newAddHeaderTransport(T http.RoundTripper) *addHeaderTransport {
	if T == nil {
		T = http.DefaultTransport
	}
	return &addHeaderTransport{T}
}

// Get HTTP call
func Get(url string) []byte {
	client := http.Client{Transport: newAddHeaderTransport(nil)}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
