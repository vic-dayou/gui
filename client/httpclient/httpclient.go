package httpclient

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client *http.Client

const xWWWFormUrlEncoded = "application/x-www-form-urlencoded"

func init() {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client = &http.Client{Timeout: 60 * time.Second, Transport: t}
}

func Post(values url.Values, link string) ([]byte, error) {
	response, err := client.Post(link, xWWWFormUrlEncoded, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {

	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(body)
	}
	return body, nil
}

func Get() {

}
